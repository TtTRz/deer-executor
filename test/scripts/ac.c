#include <stdio.h>

int main(int argc, char **argv)
{
	int a, b;
	while (~scanf("%d%d", &a, &b)) {
	    printf("%d\t\n", a+b);
	}
}
