<?php

declare(strict_types=1);

namespace DoctrineMigrations;

use Doctrine\DBAL\Schema\Schema;
use Doctrine\Migrations\AbstractMigration;

/**
 * Auto-generated Migration: Please modify to your needs!
 */
final class Version20220927090815 extends AbstractMigration
{
    public function getDescription(): string
    {
        return '';
    }

    public function up(Schema $schema): void
    {
        // this up() migration is auto-generated, please modify it to your needs
        $this->addSql('CREATE TABLE recurring_transaction (id INT AUTO_INCREMENT NOT NULL, category_id INT NOT NULL, user_id INT NOT NULL, amount NUMERIC(19, 4) NOT NULL, end_date DATE DEFAULT NULL, name VARCHAR(255) NOT NULL, period_num INT NOT NULL, period_type VARCHAR(5) NOT NULL, start_date DATE NOT NULL, summary VARCHAR(255) NOT NULL, INDEX IDX_D3509AA612469DE2 (category_id), INDEX IDX_D3509AA6A76ED395 (user_id), PRIMARY KEY(id)) DEFAULT CHARACTER SET utf8mb4 COLLATE `utf8mb4_unicode_ci` ENGINE = InnoDB');
        $this->addSql('ALTER TABLE recurring_transaction ADD CONSTRAINT FK_D3509AA612469DE2 FOREIGN KEY (category_id) REFERENCES category (id)');
        $this->addSql('ALTER TABLE recurring_transaction ADD CONSTRAINT FK_D3509AA6A76ED395 FOREIGN KEY (user_id) REFERENCES user (id)');
    }

    public function down(Schema $schema): void
    {
        // this down() migration is auto-generated, please modify it to your needs
        $this->addSql('ALTER TABLE recurring_transaction DROP FOREIGN KEY FK_D3509AA612469DE2');
        $this->addSql('ALTER TABLE recurring_transaction DROP FOREIGN KEY FK_D3509AA6A76ED395');
        $this->addSql('DROP TABLE recurring_transaction');
    }
}
