/*#include<stdio.h>
int count(char a[100]){
	for (int i = 0, int s = 0;;s++) {
		if (a[s] != '\0') {
			i++;
		}
		else {
			return i;
		}
	}
}
int main() {
	char a[100];
	gets(a);
	printf("%d",count(a));
	return 0;
}
*/

#include <time.h>
#include <stdio.h>
#include <stdlib.h>
#define N 100000
int main() {
	int a[N];
	int x;
	double time;
	double startime, endtime;
	printf("随机生成整数\n");
	for (x = 0;x <= N;x++) {
		a[x] = rand();
		printf("%d\t", a[x]);
	}
	printf("\n");

	startime = clock();

	int j, k, t;
	bool change;
	/*for(j=N-1,change=true;j>1&&change;--j){
		change=false;
	for(k=0;k<j;++k){
		if(a[k]>a[k+1]){
			t=a[k];
			a[k]=a[k+1];
			a[k+1]=j;
			change=true;
		}
	}
	}*/
	for (j = 0;j < N;j++) {
		for (k = j + 1;k < N;k++) {
			if (a[j] > a[k]) {
				t = a[k];
				a[k] = a[j];
				a[j] = t;
			}
		}
	}

	endtime = clock();

	printf("输出排序好的整数\n");
	for (x = 0;x <= N;x++) {
		printf("%d\t", a[x]);
	}
	time = (double)(endtime - startime) / CLOCKS_PER_SEC;
	printf("\n执行时间%f秒\n", time);
	return 0;
}