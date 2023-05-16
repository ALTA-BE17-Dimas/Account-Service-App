CREATE DATABASE `Account_Service_DB`;

USE `Account_Service_DB`;

CREATE TABLE `Users` (
    `id` INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `full_name` VARCHAR(255) NOT NULL,
    `identity_number` VARCHAR(50) NOT NULL,
    `birth_date` DATE NOT NULL,
    `address` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(20) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `balance` DECIMAL NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    
    CONSTRAINT `uq_Users_email` UNIQUE (`email`),
    CONSTRAINT `uq_Users_phone` UNIQUE (`phone`),
    CONSTRAINT `uq_Users_identity_number` UNIQUE (`identity_number`)
);

CREATE TABLE `Transfer_Histories`(
    `user_id_sender` INT NOT NULL,
    `user_id_recipient` INT NOT NULL,
    `amount` DECIMAL NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_Transfer_Histories_Users_Sender
        FOREIGN KEY (`user_id_sender`) REFERENCES Users(`id`),
    CONSTRAINT fk_TransferHistories_Users_Recipient
        FOREIGN KEY (`user_id_recipient`) REFERENCES Users(`id`)
);

CREATE TABLE `Top_Up_Histories`(
    `user_id` INT NOT NULL,
    `amount` DECIMAL NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_Top_Up_Histories_Users
        FOREIGN KEY (`user_id`) REFERENCES Users(`id`)
);


































