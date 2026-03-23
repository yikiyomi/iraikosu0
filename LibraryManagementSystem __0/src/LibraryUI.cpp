#include "LibraryUI.hpp"
#include<limits>
using std::cout;
using std::endl;
LibraryUI::LibraryUI(Database& database) : db(database), 
borrowMgr(database), reportGen(database), currentAdminId(0) {}

void LibraryUI::mainMenu(){
    int choice;
    do{
        cout<<"====图书管理系统====\n";
        cout<<"1.图书管理\n";
        cout<<"2.借还书\n";
        cout<<"3.统计报案\n";
        cout<<"4.系统管理\n";
        cout<<"0.退出\n";
        cout<<"输入数字以选择";
        cin>>choice;

        swich(choice){
            case 0:cout<<"感谢使用"endl;
            break;
            case 1: bookManagementMenu();break;
            case 2: readerManagementMenu(); break;
            case 3: borrowReturnMenu(); break;
            case 4: reportMenu(); break;
            case 5: systemMenu(); break;
            default:cout<<"无效"<<endl;
        }
    }while(choice!=0);
}

void LibraryUI::bookManagementMenu(){
    int choice;
    do{
        cout<<"\n---图书管理---\n";
        cout<<"1.添加图书\n";
        cout<<"2.修改图书\n";
        cout<<"3.删除图书\n";
        cout<<"4.查询图书\n";
        cout<<"5.显示所有图书\n";
        cout<<"0.返回上级\n";
        cout<<"请选择\n";
        cin>>choice;

        if(choice==1){
            Book book;
            cout<<"请输入isbn\n";book.setIsbn(getstringinput(""));
            cout<<"请输入书名\n";book.setTitle(getstringinput(""));
            cout<<"请输入作者\n";book.setAuthor(getstringinput(""));
            cout<<"请输入价格\n";book.setPrice(getintinput(""));
            cout<<"请输入总册数\n";book.setTotalCopies(getintinput(""));
            if(book.addToDatabase(db)){
                cout<<"添加成功"<<endl;
            }else{
                cout<<"添加失败"<<endl;
            }
        }
        if(choice==2){
            int id;
            cout<<"请输入 图书id"<<endl;
            cin>>id;
            auto book=Book::fetchFromDatabase(db,id);
            if(!book){
                cout<<"不存在"<<endl;
                continue;
            }
            //修改书名
            cout<<"原书名："<<book->getTitle()<<"新书名：";
            std::string newTitle=getstringinput("");
            if(!newTitle.empty())book->setTitle(newTitle);
            //作者
            cout<<"原作者："<<book->getAuthor()<<"新作者：";
            std::string newAuthor=getstringinput("");
            if(!newAuthor.empty())book->setAuthor(newAuthor);
            //价格
            cout<<"原价："<<book->getPrice()<<"新价格";
            std::int newPrice=getintinput("");
            if(!newPrice.empty())book->setPrice(newPrice);
            //册数
            cout<<"原册数："<<book->getTotalCopies()<<"新册数";
            std::int newTotalCopies=getTotalCopies("");
            if(!newTotalCopies.empty())book->setTotalCopies(newTotalCopies);
            if(book->updateInDatabase(db)){
                cout<<"修改成功"<<endl;
            }else{
                cout<<"修改失败"<<endl;
            }
        }
        if(choice==3){
            int id;
            cout<<"请输入图书id"<<endl;
            cin>>id;
            auto book=Book::fetchFromDatabase(db,id);
            if(book){
                cout<<"图书不存在"<<endl;
                continue;
            }
            cout<<"确定删除《"<<book->getTitle()<<"》吗"<<endl;
            std::string confirm=getstringinput("");
            if(confirm=="y"||confirm=="Y"){
                if(book->deleteFromDatabase(db)){
                    cout<<"删除成功"<<endl;
                }else{
                    cout<<"删除失败"<<endl;
                }
            }
        }
        if(choice==4){
            int opt;
            cout<<"按书名查询（1）按作者查询（2）"
            cin>>opt;
            std::string keyword;
            std::vector<Book> books;
            if(opt==1){
                std::cout<<"请输入书名关键字"
               keyword=getstringinput("");
               books=Book::searchByTitle(db,keyword);
            }else if(opt==2){
                std::cout<<"请输入作者关键梓"；
                keyword=getstringinput("");
                books=Book::searchByAuthor(db,keyword);
            }
            displayBooks(books);
        }
        if(choice==5){
             auto res=db.executeQuery(select * from Books limit 100);
             std::vector<Book> book;
             if(res){
                while (res->next()){
                    Book b;
                    b.setBookId(res->getInt("book_id"));
                    b.setIsbn(res->getString("isbn"));
                    b.setTitle(res->getstring("title"));
                    b.setAuthor(res->getstring("author"));
                    books,push_back(b);
                }
             }
             displayBooks(books);
        }
    }while(choice!=0);
}

void LibraryUI::readerManagementMenu(){
    //读者管理功能
    cout<<"开发中...."endl;
}

void LibraryUI::borrowReturnMenu(){
    int choice;
    do{
        cout<<"\n--借还菜单--\n";
        cout<<"1.借书\n";
        cout<<"2.还书\n";
        cout<<"3.续借\n";
        cout<<"4.发送催还提醒\n";
        cout<<"0.返回\n";
        cout<<"请选择\n";
        cin>>choice;
        
        if(choice==1){
            int readerId=getIntinput("请输入读者ID：");
            std::string isbn=getstringinput("请输入图书isbn：");
            if(borrowMgr.borrowBook(reader,isbn)){
                cout<<"借书成功"<<endl;
            }
        }
        id(choice==2){
            int recordId=getintinput("请输入借阅记录id：");
            std::cout<<"输入当前图书状况（正常|损坏|丢失）";
            std::string cond=getstringinput("");
            borrowMgr.returnBook(record,cond);
        }
        if(choice==3){
            int recordId=getintinput("请输入借阅记录id：");
            if(borrowMgr.rneewBook(record)){
                cout<<"续借成功"<<endl;
            }
        }
        if(choice==4){
            borrowMgr.sendDueDateReminders();
            borrowMgr.sendOverdueNotices();
        }
    }while(choice!=0);    
}

void LIbraryUI::reportMenu(){
    int choice;
    do{
        std::cout << "\n--- 统计报表 ---\n";
        std::cout << "1. 借阅统计\n";
        std::cout << "2. 热门图书\n";
        std::cout << "3. 读者活跃度\n";
        std::cout << "4. 逾期分析\n";
        std::cout << "5. 馆藏统计\n";
        std::cout << "0. 返回\n";
        std::cout << "请选择: ";
        choice = getintinput("");

        if(choice==1){
            DateRange range;
            cout<<"起始日期(****-**-**)";range.start=getstringinput("");
            cout<<"终止日期"；range.end=getstringinput("");
            reportGen.generateBorrowStatistics(range);
        }
        if(choice==2){
            int top=getintinput("请输入前n本");
            reportGen.generatePopularBookReport(top);
        }
        if(choice==3){
            reportGen.generateReaderActivityReport();
        }
        if(choice==4){
            reportGen.generateOverdueAnalysis();
        }
        if(choice==5){
            reportGen.generateInventoryReport();
        }
    }while(choice!=0);
}

void LibraryUI::systemMenu(){
    cout<<"功能开发中"<<endl;
}

void LibraryUI::displayBooks(const std::vector<Book>& books){
    if(books.empty()){
        cout<<"没找到"<<endl;
        return ;
    }
    
    for(const auto& b:books){
    cout<<"ID："<<b.getBookId()<< ", ISBN: " << b.getIsbn()
    << ", 书名: " << b.getTitle() << ", 作者: " << b.getAuthor()
    << ", 状态: " << b.getStatus() << std::endl;
    }
}

void LibraryUI::displayReaders(const std::vector<Reader>& readers) {
    // 类似实现...
}

void LibraryUI::displayBorrowRecords(const std::vector<BorrowRecord>& records) {
    // 类似实现...
}

int LibraryUI::getintinput(const std::string&prompt){
    int value;
    std::cin>>value;
    std::cin.ignore(std::numeric_limits<std::streamsize>::max(),'\n');
    return value;
}

std::string LibraryUI::getstring(const std::string& prompt){
    std::string line;
    std::getline(std::cin,line);
    return line;
}

std::string LibraryUI::getdateinput(const std::string&prompt){
    //验证格式说是
    return getstringinput(prompt);
}