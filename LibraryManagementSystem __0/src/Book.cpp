#include "Book.hpp"
#include <sstream>

Book::Book() : bookId(0), price(0.0), totalCopies(1), availableCopies(1) {}

Book::Book(const std::string& isbn, const std::string& title, const std::string& author)
    : bookId(0), isbn(isbn), title(title), author(author), price(0.0),
      totalCopies(1), availableCopies(1), status("可借") {}

bool Book::addToDatabase(Database& db) {
    // 使用 ? 占位符，禁止字符串拼接
    std::string sql = "INSERT INTO Books (isbn, title, author, publisher, "
                      "publish_date, category, price, total_copies, "
                      "available_copies, location, status) "
                      "VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)";
    
    std::vector<std::string> params = {
        isbn,
        title,
        author,
        publisher,
        publishDate,
        category,
        std::to_string(price),        // 数值转字符串
        std::to_string(totalCopies),
        std::to_string(availableCopies),
        location,
        status
    };
    
    return db.executeUpdate(sql, params) > 0;  // 
}

bool Book::updateInDatabase(Database&db){
    std::string sql="update Books set isbn=?,title=?,author=?,
    publisher=?,publish_date=?,category=?,price=?,total_copies=?,
    available_copice=?,location=?,status=?,WHERE book_id=?";
     std::vector<std::string> params = {
        isbn,                    // 1. isbn
        title,                   // 2. title
        author,                  // 3. author
        publisher,               // 4. publisher
        publishDate,             // 5. publish_date
        category,                // 6. category
        std::to_string(price),   // 7. price
        std::to_string(totalCopies),      // 8. total_copies
        std::to_string(availableCopies), // 9. available_copies
        location,                // 10. location
        status,                  // 11. status
        std::to_string(bookId)   // 12. WHERE book_id=?
    };
    
    return db.executeUpdate(sql, params) > 0;
}

bool Book::deleteFromDatabase(Database& db) {
    std::string sql = "DELETE FROM Books WHERE book_id=" + std::to_string(bookId);
    return db.executeUpdate(sql) > 0;
}

std::unique_ptr<Book> Book::fetchFromDatabase(Database& db, int bookId) {
    std::string sql = "SELECT * FROM Books WHERE book_id=" + std::to_string(bookId);
    auto res = db.executeQuery(sql);
    if (res && res->next()) {
        auto book = std::make_unique<Book>();
        book->setBookId(res->getInt("book_id"));
        book->setIsbn(res->getString("isbn"));
        book->setTitle(res->getString("title"));
        book->setAuthor(res->getString("author"));
        book->setPublisher(res->getString("publisher"));
        book->setPublishDate(res->getString("publish_date"));
        book->setCategory(res->getString("category"));
        book->setPrice(res->getDouble("price"));
        book->setTotalCopies(res->getInt("total_copies"));
        book->setAvailableCopies(res->getInt("available_copies"));
        book->setLocation(res->getString("location"));
        book->setStatus(res->getString("status"));
        return book;
    }
    return nullptr;
}

std::vector<Book> Book::searchByTitle(Database& db, const std::string& keyword) {
    std::vector<Book> books;
    std::string sql = "SELECT * FROM Books WHERE title LIKE '%" + keyword + "%'";
    auto res = db.executeQuery(sql);
    if (res) {
        while (res->next()) {
            Book book;
            book.setBookId(res->getInt("book_id"));
            book.setIsbn(res->getString("isbn"));
            book.setTitle(res->getString("title"));
            book.setAuthor(res->getString("author"));
            book->setPublisher(res->getString("publisher"));
            book.setPublishDate(res->getString("publish_date"));
            book.setCategory(res->getString("category"));
            book.setPrice(res->getDouble("price"));
            book.setTotalCopies(res->getInt("total_copies"));
            book.setAvailableCopies(res->getInt("available_copies"));
            book.setLocation(res->getString("location"));
            book.setStatus(res->getString("status"));
            books.push_back(book);
        }
    }
    return books;
}

std::vector<Book> Book::searchByAuthor(Database& db, const std::string& author) {
    std::vector<Book> books;
    std::string sql = "SELECT * FROM Books WHERE author LIKE '%" + author + "%'";
    auto res = db.executeQuery(sql);
    if (res) {
        while (res->next()) {
            Book book;
            book.setBookId(res->getInt("book_id"));
            book.setIsbn(res->getString("isbn"));
            book.setTitle(res->getString("title"));
            book.setAuthor(res->getString("author"));
            book->setPublisher(res->getString("publisher"));
            book.setPublishDate(res->getString("publish_date"));
            book.setCategory(res->getString("category"));
            book.setPrice(res->getDouble("price"));
            book.setTotalCopies(res->getInt("total_copies"));
            book.setAvailableCopies(res->getInt("available_copies"));
            book.setLocation(res->getString("location"));
            book.setStatus(res->getString("status"));
            books.push_back(book);
        }
    }
    return books;
}

bool Book::decreaseAvailableCopies() {
    if (availableCopies > 0) {
        availableCopies--;
        return true;
    }
    return false;
}

bool Book::increaseAvailableCopies() {
    availableCopies++;
    return true;
}