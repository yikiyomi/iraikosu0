#include<iostream>
#include<string>
#define MAX 1000

using namespace std;

struct Person {
	string m_Name;//名字
	int m_Sex;//性别
	int m_Age;//年龄
	string m_Phone;//手机号码
	string m_Addr;//住址
};

struct Addressbooks {
	struct Person personArray[MAX];
	int m_Size;
};
//菜单
void showMenu() {
	cout << "1.添加联系人" << endl;
	cout << "2.显示联系人" << endl;
	cout << "3.删除联系人" << endl;
	cout << "4.查找联系人" << endl;
	cout << "5.修改联系人" << endl;
	cout << "6.清空联系人" << endl;
	cout << "0.退出通讯录" << endl;
}
//添加联系人
void addPerson(Addressbooks * abs) {
	string name;
	cout << "请输入姓名" << endl;
	cin >> name;
	abs->personArray[abs->m_Size].m_Name = name;

	int sex;
	cout << "请输入性别（男1女2）" << endl;
	while (true) {
		cin >> sex;
		if (sex==1 || sex==2) {
		abs->personArray[abs->m_Size].m_Sex = sex;
		break;
		}
		else {
			cout << "输入有误请重新输入" << endl;
		}
	}
	
	

	int age;
	cout << "请输入年龄" << endl;
	cin >> age;
	abs->personArray[abs->m_Size].m_Age = age;

	string phone;
	cout << "请输入手机号" << endl;
	cin >> phone;
	abs->personArray[abs->m_Size].m_Phone = phone;

	string address;
	cout << "请输入家庭住址" << endl;
	cin >> address;
	abs->personArray[abs->m_Size].m_Addr = address;

	abs->m_Size++;
	cout << "添加成功" << endl;

	system("pause");
	system("cls");
}
//显示联系人
void shouPerson(Addressbooks* abs) {
	if (abs->m_Size == 0) {
		cout << "联系人为空" << endl;
	}
	else {
		for (int i = 0;i < abs->m_Size;i++) {
			cout << "姓名为" << abs->personArray[i].m_Name << endl;
			cout << "年龄为" << abs->personArray[i].m_Age << endl;

			if (abs->personArray[i].m_Sex == 1) {
				cout << "性别为" << "男" << endl;
			}
			else {
				cout << "性别为" << "女" << endl;
			}
			
			cout << "手机号码是：" << abs->personArray[i].m_Phone << endl;
			cout << "家庭住址在：" << abs->personArray[i].m_Addr << endl;
			cout << "*******" << endl;
		}
	}
	system("pause");
	system("cls");
}
//检测联系人是否存在
int isExist(Addressbooks* abs, string name) {
	for (int t = 0;t < abs->m_Size;t++) {
		if (abs->personArray[t].m_Name == name) {
			return t;
		}
	}
	return -1;
}
//查找联系人
void findPerson(Addressbooks* abs) {
	cout << "请输入需查找的联系人名字" << endl;
	string name;
	cin >> name;
	if (isExist(abs, name) == -1) {
		cout << "该联系人不存在" << endl;
	}
	else {
		cout << "姓名" <<abs->personArray[isExist(abs, name)].m_Name<< endl;
		cout << "性别" << abs->personArray[isExist(abs, name)].m_Sex<<endl;
		cout << "年龄" << abs->personArray[isExist(abs, name)].m_Age<<endl;
		cout << "电话" << abs->personArray[isExist(abs, name)].m_Phone<< endl;
		cout << "家庭住址" << abs->personArray[isExist(abs, name)].m_Addr<< endl;
	}
	system("pause");
	system("cls");
}
//修改联系人信息
void modifyPerson(Addressbooks * abs) {
	cout << "请输入需要修改的联系人名字" << endl;
	string name;
	cin >> name;
	int i = isExist(abs, name);
	if (i==-1) {
		cout << "联系人不存在" << endl;
	}
	else {
		string name;
		cout << "请输入姓名" << endl;
		cin >> name;
		abs->personArray[i].m_Name = name;

		int sex;
	     cout << "请输入性别（男1女2）" << endl;
	     while (true) {
		cin >> sex;
		if (sex==1 || sex==2) {
		abs->personArray[i].m_Sex = sex;
		break;
		}
		else {
			cout << "输入有误请重新输入" << endl;
		}
	    }

		 int age;
		 cout << "请输入年龄" << endl;
		 cin >> age;
		 abs->personArray[i].m_Age = age;

		 string phone;
		 cout << "请输入手机号" << endl;
		 cin >> phone;
		 abs->personArray[i].m_Phone = phone;

		 string address;
		 cout << "请输入家庭住址" << endl;
		 cin >> address;
		 abs->personArray[i].m_Addr = address;

		 cout << "修改成功" << endl;
	}
	system("pause");
	system("cls");
}
//清空
void cleanPerson(Addressbooks * abs) {
	abs->m_Size = 0;
	cout << "通讯录已清空" << endl;
	system("pause");
	system("cls");
}
int main() {
	Addressbooks abs;//创建通讯录结构体变量
	abs.m_Size = 0;//人员个数

	int select = 0;//输入变量
	while (true) {
		//调用菜单
		showMenu();

		cin >> select;

		switch (select)
		{
		case 1:
			addPerson(&abs);
			break;
		case 2:
			shouPerson(&abs);
			break;
		case 3:
		{
			cout << "请输入需删除联系人的名字：" << endl;
			string name;
			cin >> name;
			if (isExist(&abs, name) == -1) {
				cout << "该联系人不存在" << endl;
			}
			else {
				for (int i = isExist(&abs, name);i<abs.m_Size;i++) {
					abs.personArray[i]= abs.personArray[i + 1];
				}
				abs.m_Size--;
				cout << "删除成功" << endl;
			}
		}
		system("pause");
		system("cls");
		break;
		case 4:
			findPerson(&abs);
			break;
		case 5:
			modifyPerson(&abs);
			break;
		case 6:
			cleanPerson(&abs);
			break;
		case 0:
			cout << "欢迎下次使用" << endl;
			system("pause");
			return 0;
			break;
		default:
			cout << "请输入正确的数字" << endl;
			break;
		}
	}

	system("pause");
	return 0;
}