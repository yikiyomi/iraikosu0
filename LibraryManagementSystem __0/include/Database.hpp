#ifndef DATABASE_HPP
#define DATABASE_HPP

#include <mysql_driver.h>
#include <mysql_connection.h>
#include <cppconn/statement.h>
#include <cppconn/resultset.h>
#include <cppconn/exception.h>
#include <cppconn/prepared_statement.h>

#include <string>
#include <memory>
#include <iostream>

class Database{
    private:
    std::unique_ptr<sql::Connection>conn;
    std::string host,user,pass,database;

    public:
    Database(conststd::string& host,cost std::string&user,
            const std::string& pass, const std::string& db);
            ~Database();

            bool connect();
            void disconnect();

            // 执行查询，返回结果集（调用者负责释放）
            std::unique_ptr<sql::ResultSet> executeQuery(const std::string& sql);
            // 执行更新，返回影响行数
            int executeUpdate(const std::string& sql);

            // 事务支持
            void beginTransaction();
            void commit();
            void rollback();

            // 获取预处理语句（防止SQL注入）
            std::unique_ptr<sql::PreparedStatement> prepareStatement(const std::string& sql);
};

#endif // DATABASE_HPP
