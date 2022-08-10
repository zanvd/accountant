<?php

declare(strict_types=1);

namespace DoctrineMigrations;

use Doctrine\DBAL\Schema\Schema;
use Doctrine\Migrations\AbstractMigration;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
final class Version20220810150300 extends AbstractMigration
{
    public function getDescription(): string
    {
        return '';
    }

    public function up(Schema $schema): void
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->addSql('CREATE TABLE category (id INT AUTO_INCREMENT NOT NULL, user_id INT NOT NULL, color VARCHAR(7) NOT NULL, description VARCHAR(150) NOT NULL, name VARCHAR(30) NOT NULL, text_color VARCHAR(7) NOT NULL, INDEX IDX_64C19C1A76ED395 (user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('CREATE TABLE transaction (id INT AUTO_INCREMENT NOT NULL, category_id INT NOT NULL, user_id INT NOT NULL, amount NUMERIC(19, 4) NOT NULL, name VARCHAR(255) NOT NULL, summary VARCHAR(255) NOT NULL, transaction_date DATE NOT NULL, INDEX IDX_723705D112469DE2 (category_id), INDEX IDX_723705D1A76ED395 (user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('CREATE TABLE transaction_template (id INT AUTO_INCREMENT NOT NULL, category_id INT DEFAULT NULL, user_id INT NOT NULL, name VARCHAR(255) NOT NULL, position INT NOT NULL, transaction_type VARCHAR(7) NOT NULL, INDEX IDX_D91236EB12469DE2 (category_id), INDEX IDX_D91236EBA76ED395 (user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('CREATE TABLE user (id INT AUTO_INCREMENT NOT NULL, email VARCHAR(255) NOT NULL, is_verified TINYINT(1) NOT NULL, password VARCHAR(255) NOT NULL, roles JSON NOT NULL, username VARCHAR(255) NOT NULL, UNIQUE INDEX UNIQ_8D93D649F85E0677 (username), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('ALTER TABLE category ADD CONSTRAINT FK_64C19C1A76ED395 FOREIGN KEY (user_id) REFERENCES user (id)');
        $this->addSql('ALTER TABLE transaction ADD CONSTRAINT FK_723705D112469DE2 FOREIGN KEY (category_id) REFERENCES category (id)');
        $this->addSql('ALTER TABLE transaction ADD CONSTRAINT FK_723705D1A76ED395 FOREIGN KEY (user_id) REFERENCES user (id)');
        $this->addSql('ALTER TABLE transaction_template ADD CONSTRAINT FK_D91236EB12469DE2 FOREIGN KEY (category_id) REFERENCES category (id)');
        $this->addSql('ALTER TABLE transaction_template ADD CONSTRAINT FK_D91236EBA76ED395 FOREIGN KEY (user_id) REFERENCES user (id)');
    }

    public function down(Schema $schema): void
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->addSql('ALTER TABLE transaction DROP FOREIGN KEY FK_723705D112469DE2');
        $this->addSql('ALTER TABLE transaction_template DROP FOREIGN KEY FK_D91236EB12469DE2');
        $this->addSql('ALTER TABLE category DROP FOREIGN KEY FK_64C19C1A76ED395');
        $this->addSql('ALTER TABLE transaction DROP FOREIGN KEY FK_723705D1A76ED395');
        $this->addSql('ALTER TABLE transaction_template DROP FOREIGN KEY FK_D91236EBA76ED395');
        $this->addSql('DROP TABLE category');
        $this->addSql('DROP TABLE transaction');
        $this->addSql('DROP TABLE transaction_template');
        $this->addSql('DROP TABLE user');
    }
}
