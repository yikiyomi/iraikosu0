#ifndef READER_HPP
#define READER_HPP

#include "Database.hpp"
#include <string>
#include <vector>
#include <memory>

class Reader {
private:
    int readerId;
    std::string cardNumber;
    std::string name;
    std::string gender;
    std::string phone;
    std::string email;
    std::string address;
    std::string readerType;
    int maxBorrowLimit;
    int currentBorrowed;
    std::string registrationDate;
    std::string expiryDate;
    std::string status;
    double penalty;

public:
    Reader();
    Reader(const std::string& cardNo, const std::string& name);

    // Getters/Setters...
    int getReaderId() const { return readerId; }
    void setReaderId(int id) { readerId = id; }
    std::string getCardNumber() const { return cardNumber; }
    void setCardNumber(const std::string& cn) { cardNumber = cn; }
    std::string getName() const { return name; }
    void setName(const std::string& n) { name = n; }
    std::string getGender() const { return gender; }
    void setGender(const std::string& g) { gender = g; }
    std::string getPhone() const { return phone; }
    void setPhone(const std::string& p) { phone = p; }
    std::string getEmail() const { return email; }
    void setEmail(const std::string& e) { email = e; }
    std::string getAddress() const { return address; }
    void setAddress(const std::string& a) { address = a; }
    std::string getReaderType() const { return readerType; }
    void setReaderType(const std::string& t) { readerType = t; }
    int getMaxBorrowLimit() const { return maxBorrowLimit; }
    void setMaxBorrowLimit(int limit) { maxBorrowLimit = limit; }
    int getCurrentBorrowed() const { return currentBorrowed; }
    void setCurrentBorrowed(int n) { currentBorrowed = n; }
    std::string getRegistrationDate() const { return registrationDate; }
    void setRegistrationDate(const std::string& d) { registrationDate = d; }
    std::string getExpiryDate() const { return expiryDate; }
    void setExpiryDate(const std::string& d) { expiryDate = d; }
    std::string getStatus() const { return status; }
    void setStatus(const std::string& s) { status = s; }
    double getPenalty() const { return penalty; }
    void setPenalty(double p) { penalty = p; }

    // 注册和管理
    bool registerReader(Database& db);
    bool updateInfo(Database& db);
    bool deleteReader(Database& db);
    static std::unique_ptr<Reader> fetchFromDatabase(Database& db, int readerId);
    static std::unique_ptr<Reader> fetchByCardNumber(Database& db, const std::string& cardNo);

    // 借阅功能
    bool canBorrowMore() const { return currentBorrowed < maxBorrowLimit && status == "正常"; }
    bool hasOverdueBooks(Database& db);
    bool increaseBorrowedCount();
    bool decreaseBorrowedCount();
    void addPenalty(double amount) { penalty += amount; }
};

#endif // READER_HPP