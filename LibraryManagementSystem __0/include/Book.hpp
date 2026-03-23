#ifndef BOOK_HPP
#define BOOK_HPP

#include "Database.hpp"
#include <string>
#include <vector>
#include <memory>

class Book{
    private:
    int bookId;
    std::string isbn;
    std::string title;
    std::string author;
    std::string publisher;
    std::string publishDate;
    std::string category;
    double price;
    int totalCopies;
    int availableCopies;
    std::string location;
    std::string status;

    public:
    Book();
    Book(const std::string& isbn,const std::string&title,
        const std::string::&author);


        //getters/setters
        int getBookId() const {return bookId;}
        void setBookId(int id){
            bookId=id;
        }
        std::string getIsbn() const{
            return isbn;
        }
        void setIsbn(const std:: string &s){
            isbn=s;
        }
        std::string getTitle() const { return title; }
        void setTitle(const std::string& t) { title = t; }
        std::string getAuthor() const { return author; }
        void setAuthor(const std::string& a) { author = a; }
        std::string getPublisher() const { return publisher; }
        void setPublisher(const std::string& p) { publisher = p; }
        std::string getPublishDate() const { return publishDate; }
        void setPublishDate(const std::string& d) { publishDate = d; }
        std::string getCategory() const { return category; }
        void setCategory(const std::string& c) { category = c; }
        double getPrice() const { return price; }
        void setPrice(double p) { price = p; }
        int getTotalCopies() const { return totalCopies; }
        void setTotalCopies(int n) { totalCopies = n; }
        int getAvailableCopies() const { return availableCopies; }
        void setAvailableCopies(int n) { availableCopies = n; }
        std::string getLocation() const { return location; }
        void setLocation(const std::string& loc) { location = loc; }
        std::string getStatus() const { return status; }
        void setStatus(const std::string& s) { status = s; }

        // CRUD 操作
        bool addToDatabase(Database& db);
        bool updateInDatabase(Database& db);
        bool deleteFromDatabase(Database& db);
        static std::unique_ptr<Book> fetchFromDatabase(Database& db, int bookId);
        static std::vector<Book> searchByTitle(Database& db, const std::string& keyword);
        static std::vector<Book> searchByAuthor(Database& db, const std::string& author);

        // 借阅相关
        bool isAvailable() const { return availableCopies > 0 && status == "可借"; }
        bool decreaseAvailableCopies();
        bool increaseAvailableCopies();
};

#endif // BOOK_HPP