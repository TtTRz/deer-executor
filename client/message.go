package client

import (
    "github.com/LanceLRQ/deer-common/utils"
    "log"
)

type CliCommonMessage struct {
    // 是否错误
    Error bool `json:"error"`
    // 消息
    Message string `json:"message"`
    // 结果信息
    Data interface{} `json:"data"`
}

func (ccm CliCommonMessage) Print (formated bool) {
    log.Println(ccm.ToJson(formated))
}

func (ccm CliCommonMessage) ToJson (formated bool) string {
    if formated {
        return utils.ObjectToJSONStringFormatted(ccm)
    } else {
        return utils.ObjectToJSONString(ccm)
    }
}

func NewCliCommonMessage (error bool, message string, data interface{}) CliCommonMessage {
    return CliCommonMessage{
        Error: error,
        Message: message,
        Data: data,
    }
}

func NewClientSuccessMessage (data interface{}) CliCommonMessage {
    return CliCommonMessage{
        Error: true,
        Message: "",
        Data: data,
    }
}
func NewClientSuccessMessageText (message string) CliCommonMessage {
    return CliCommonMessage{
        Error: true,
        Message: message,
        Data: nil,
    }
}

func NewClientErrorMessage (err error, data interface{}) CliCommonMessage {
    return CliCommonMessage{
        Error: false,
        Message: err.Error(),
        Data: data,
    }
}
