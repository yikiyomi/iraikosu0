#include "Database.hpp"

Database::Database(const std::string& host, const std::string& user,
                   const std::string& pass, const std::string& db)
    : host(host), user(user), pass(pass), database(db), conn(nullptr) {}

Database::~Database(){
    disconnect();
}

bool Database::connect(){
    try{
        sql::mysql::Mysql_Driver*driver=sql::mysql::get_mysql_driver_instance();
        conn=std::unique_ptr<sql::Connection>(driver->connect(host, user, pass));
    conn->setSchema(database);
    std::cout<<"数据库连接成功"<<std::endl;
    return ture;
    }catch(sql::SQLException&e){
        std::cerr<<"连接失败"<<e.what()<<std::endl;
        return false;
    }

}

void Database::disconnect(){
    if(conn){
        conn->close();
        conn.reset();
    }
}

std::unique_ptr<sql::ResultSet> Database::executeQuery(const std::string& sql) {
    try {
        std::unique_ptr<sql::Statement> stmt(conn->createStatement());
        std::unique_ptr<sql::ResultSet> res(stmt->executeQuery(sql));
        return res;
    } catch (sql::SQLException& e) {
        std::cerr << "查询执行失败: " << e.what() << std::endl;
        return nullptr;
    }
}

int Database::executeUpdate(const std::string& sql) {
    try {
        std::unique_ptr<sql::Statement> stmt(conn->createStatement());
        return stmt->executeUpdate(sql);
    } catch (sql::SQLException& e) {
        std::cerr << "更新执行失败: " << e.what() << std::endl;
        return -1;
    }
}
void Database::beginTransaction() {
    executeUpdate("START TRANSACTION");
}

void Database::commit() {
    executeUpdate("COMMIT");
}

void Database::rollback() {
    executeUpdate("ROLLBACK");
}

std::unique_ptr<sql::PreparedStatement> Database::prepareStatement(const std::string& sql) {
    try {
        return std::unique_ptr<sql::PreparedStatement>(conn->prepareStatement(sql));
    } catch (sql::SQLException& e) {
        std::cerr << "预处理语句失败: " << e.what() << std::endl;
        return nullptr;
    }
}