class Database {
private:
sql::Connection conn;
std::string host,user,pass,database;
public:
Database(const std::string&host,const std::string&user,
         const std::string&pass,const std::string&db);
        ~Database();
        bool connect();
        void disconnect();
        //执行查询
        sql::Resultset*executeQuery(const std::string& sql);
        int executeUpdate(const std::string& sql); 
        //事务支持
        void begintransaction();
        void commit();
        void rollback();
};
       