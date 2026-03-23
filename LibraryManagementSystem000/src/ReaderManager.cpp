class ReaderManager {
public:
    // 读者注册
    bool registerNewReader(const ReaderInfo& info);
    
    // 读者证办理
    LibraryCard issueLibraryCard(int readerId);
    
    // 读者信息更新
    bool updateReaderInfo(int readerId, const ReaderInfo& newInfo);
    
    // 读者证挂失/解挂
    bool reportCardLost(int readerId);
    bool unblockCard(int readerId);
    
    // 读者注销
    bool deregisterReader(int readerId, const std::string& reason);
    
    // 批量导入读者
    bool importReadersFromCSV(const std::string& filename);
};