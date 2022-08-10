<?php

namespace App\Entity;

use App\Enum\TransactionType;
use App\Repository\TransactionTemplateRepository;
use Doctrine\ORM\Mapping as ORM;

#[ORM\Entity(repositoryClass: TransactionTemplateRepository::class)]
class TransactionTemplate
{
    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column(type: 'integer')]
    private int $id;

    #[ORM\ManyToOne(targetEntity: Category::class, inversedBy: 'transactionTemplates')]
    private Category $category;

    #[ORM\Column(type: 'string', length: 255)]
    private string $name;

    #[ORM\Column(type: 'integer')]
    private int $position;

    #[ORM\Column(type: 'string', enumType: TransactionType::class, length: 7)]
    private TransactionType $transactionType;

    #[ORM\ManyToOne(targetEntity: User::class, inversedBy: 'transactionTemplates')]
    #[ORM\JoinColumn(nullable: false)]
    private $user;

    public function getId(): ?int
    {
        return $this->id;
    }

    public function getCategory(): ?Category
    {
        return $this->category;
    }

    public function setCategory(?Category $category): self
    {
        $this->category = $category;

        return $this;
    }

    public function getName(): ?string
    {
        return $this->name;
    }

    public function setName(string $name): self
    {
        $this->name = $name;

        return $this;
    }

    public function getPosition(): ?int
    {
        return $this->position;
    }

    public function setPosition(int $position): self
    {
        $this->position = $position;

        return $this;
    }

    public function getTransactionType(): ?TransactionType
    {
        return $this->transactionType;
    }

    public function setTransactionType(TransactionType $transactionType): self
    {
        $this->transactionType = $transactionType;

        return $this;
    }

    public function isIncome(): bool
    {
        return $this->transactionType === TransactionType::Income;
    }

    public function getUser(): ?User
    {
        return $this->user;
    }

    public function setUser(?User $user): self
    {
        $this->user = $user;

        return $this;
    }
}
