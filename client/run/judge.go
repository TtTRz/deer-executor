// +build linux darwin

package run

import (
    "github.com/LanceLRQ/deer-common/logger"
    "github.com/LanceLRQ/deer-common/persistence"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/LanceLRQ/deer-executor/v2/executor"
    "github.com/pkg/errors"
    uuid "github.com/satori/go.uuid"
    "github.com/urfave/cli/v2"
    "log"
    "os"
    "path/filepath"
)

// 执行一次完整的评测
func runOnceJudge(options *RunOption) (*commonStructs.JudgeResult, *executor.JudgeSession, error) {
    // create session
    session, err := executor.NewSessionWithLog(options.ConfigFile, options.ShowLog, options.LogLevel)
    if err != nil {
        return nil, nil, err
    }
    if options.Language != "" {
        session.CodeLangName = options.Language
    }
    if session.JudgeConfig.SpecialJudge.Mode > 0 {
        // 特判时需要检查library目录
        libDir, err := filepath.Abs(options.LibraryDir)
        if err != nil {
            return nil, nil, errors.Errorf("get library root error: %s", err.Error())
        }
        if s, err := os.Stat(libDir); err != nil {
            return nil, nil, errors.Errorf("library root not exists")
        } else {
            if !s.IsDir() {
                return nil, nil, errors.Errorf("library root not a directory")
            }
        }
        session.LibraryDir = libDir
    }
    // init files
    if options.WorkDir != "" {
        workDirAbsPath, err := filepath.Abs(options.WorkDir)
        if err != nil {
            return nil, nil, err
        }
        session.ConfigDir = workDirAbsPath
        session.JudgeConfig.ConfigDir = session.ConfigDir
    }
    session.CodeFile = options.CodePath
    session.SessionId = options.SessionId
    session.SessionRoot = options.SessionRoot
    // create session info
    if session.SessionId == "" {
        session.SessionId = uuid.NewV1().String()
    }
    if session.SessionRoot == "" {
        session.SessionRoot = "/tmp"
    }
    // 初始化session dir
    sessionDir, err := utils.GetSessionDir(session.SessionRoot, session.SessionId)
    if err != nil {
        log.Fatal(err)
        return nil, nil, err
    }
    session.SessionDir = sessionDir
    // start judgement
    judgeResult := session.RunJudge()
    return &judgeResult, session, nil
}

// 执行CLI的评测
func runUserJudge (c *cli.Context, configFile, workDir string) (*commonStructs.JudgeResult, error) {
    // parse params
    persistenceOn := c.String("persistence") != ""
    digitalSign := c.Bool("sign")
    compressorType := uint8(1)
    if c.String("compress") == "none" {
        compressorType = uint8(0)
    }
    jOption := persistence.JudgeResultPersisOptions{
        CompressorType: compressorType,
        SaveAcceptedData: c.Bool("save-ac-data"),
    }
    jOption.OutFile = c.String("persistence")
    // 是否要持久化结果
    if persistenceOn {
        if digitalSign {
            if c.String("passphrase") != "" {
                log.Println("[warn] Using a password on the command line interface can be insecure.")
            }
            passphrase := []byte(c.String("passphrase"))
            pem, err := persistence.GetArmorPublicKey(c.String("gpg-key"), passphrase)
            if err != nil {
                return nil, err
            }
            jOption.DigitalSign = true
            jOption.DigitalPEM = pem
        }
    }

    isBenchmarkMode := c.Int("benchmark") > 1

    // 获取log等级
    var logLevel int
    showLog := false
    if c.Bool("log") {
        showLog = true
        var ok bool
        logLevelStr := c.String("log-level")
        logLevel, ok = logger.LogLevelStrMapping[logLevelStr]
        if !ok {
            logLevel = 0
        }
    }
    showLog = !isBenchmarkMode && showLog

    // 构建运行选项
    rOptions := &RunOption{
        Clean: !c.Bool("no-clean"),
        ShowLog: showLog,
        LogLevel: logLevel,
        WorkDir: workDir,
        ConfigFile: configFile,
        Language: c.String("language"),
        LibraryDir: c.String("library"),
        CodePath: c.Args().Get(1),
        SessionId: c.String("session-id"),
        SessionRoot: c.String("session-root"),
    }

    if persistenceOn {
        rOptions.Persistence = &jOption
    }

    // 执行评测
    _, judgeResult, err := StartJudgement(rOptions)
    if err != nil {
        return nil, err
    }

    if !c.Bool("detail") {
        judgeResult.TestCases = nil
        judgeResult.JudgeLogs = nil
    }

    return judgeResult, nil
}