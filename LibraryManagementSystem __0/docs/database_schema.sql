-- 创建数据库
CREATE DATABASE IF NOT EXISTS library DEFAULT CHARACTER SET utf8mb4;
USE library;

-- 图书信息表
CREATE TABLE Books (
    book_id INT PRIMARY KEY AUTO_INCREMENT,
    isbn VARCHAR(20) UNIQUE NOT NULL,
    title VARCHAR(200) NOT NULL,
    author VARCHAR(100) NOT NULL,
    publisher VARCHAR(100),
    publish_date DATE,
    category VARCHAR(50),
    price NUMERIC(8,2),
    total_copies INT DEFAULT 1,
    available_copies INT DEFAULT 1,
    location VARCHAR(50),
    status ENUM('可借','借出','修补中','遗失') DEFAULT '可借',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 读者信息表
CREATE TABLE Readers (
    reader_id INT PRIMARY KEY AUTO_INCREMENT,
    card_number VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    gender ENUM('男','女','其他'),
    phone VARCHAR(20),
    email VARCHAR(100),
    address VARCHAR(200),
    reader_type ENUM('学生','教师','职工','其他'),
    max_borrow_limit INT DEFAULT 5,
    current_borrowed INT DEFAULT 0,
    registration_date DATE NOT NULL,
    expiry_date DATE NOT NULL,
    status ENUM('正常','挂失','冻结','注销') DEFAULT '正常',
    penalty NUMERIC(8,2) DEFAULT 0.00
);

-- 借阅记录表
CREATE TABLE BorrowRecords (
    record_id INT PRIMARY KEY AUTO_INCREMENT,
    reader_id INT NOT NULL,
    book_id INT NOT NULL,
    borrow_date DATETIME NOT NULL,
    due_date DATE NOT NULL,
    actual_return_date DATETIME,
    renew_count INT DEFAULT 0,
    overdue_days INT DEFAULT 0,
    fine_amount NUMERIC(8,2) DEFAULT 0.00,
    status ENUM('借出','已还','逾期','丢失') DEFAULT '借出',
    FOREIGN KEY (reader_id) REFERENCES Readers(reader_id),
    FOREIGN KEY (book_id) REFERENCES Books(book_id),
    INDEX idx_reader (reader_id),
    INDEX idx_book (book_id),
    INDEX idx_due_date (due_date)
);

-- 管理员表
CREATE TABLE Administrators (
    admin_id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(50),
    role ENUM('超级管理员','图书管理员','读者管理员') DEFAULT '图书管理员',
    last_login DATETIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 预约记录表
CREATE TABLE Reservations (
    reservation_id INT PRIMARY KEY AUTO_INCREMENT,
    reader_id INT NOT NULL,
    book_id INT NOT NULL,
    reserve_date DATETIME NOT NULL,
    expiry_date DATETIME,
    status ENUM('等待中','已通知','已取消','已过期') DEFAULT '等待中',
    FOREIGN KEY (reader_id) REFERENCES Readers(reader_id),
    FOREIGN KEY (book_id) REFERENCES Books(book_id)
);

-- 罚款记录表
CREATE TABLE Fines (
    fine_id INT PRIMARY KEY AUTO_INCREMENT,
    reader_id INT NOT NULL,
    amount NUMERIC(8,2) NOT NULL,
    reason ENUM('逾期','损坏','丢失') NOT NULL,
    status ENUM('未缴','已缴','豁免') DEFAULT '未缴',
    create_date DATETIME NOT NULL,
    pay_date DATETIME,
    FOREIGN KEY (reader_id) REFERENCES Readers(reader_id)
);

-- 视图
CREATE VIEW CurrentBorrowing AS
SELECT 
    r.name AS reader_name,
    r.card_number,
    b.title,
    b.isbn,
    br.borrow_date,
    br.due_date,
    DATEDIFF(br.due_date, CURDATE()) AS remaining_days
FROM BorrowRecords br
JOIN Readers r ON br.reader_id = r.reader_id
JOIN Books b ON br.book_id = b.book_id
WHERE br.status = '借出';

CREATE VIEW OverdueBooks AS
SELECT 
    r.name AS reader_name,
    r.card_number,
    b.title,
    br.due_date,
    DATEDIFF(CURDATE(), br.due_date) AS overdue_days,
    b.price
FROM BorrowRecords br
JOIN Readers r ON br.reader_id = r.reader_id
JOIN Books b ON br.book_id = b.book_id
WHERE br.status = '借出' 
AND br.due_date < CURDATE();