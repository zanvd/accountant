<?php

namespace App\Entity;

use App\Enum\TransactionType;
use App\Repository\TransactionTemplateRepository;
use Doctrine\ORM\Mapping as ORM;
use Symfony\Component\Validator\Constraints as Assert;

#[ORM\Entity(repositoryClass: TransactionTemplateRepository::class)]
class TransactionTemplate
{
    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column(type: 'integer')]
    private int $id;

    #[ORM\ManyToOne(targetEntity: Category::class, inversedBy: 'transactionTemplates')]
    private Category $category;

    #[Assert\NotBlank(message: 'Name is required.')]
    #[ORM\Column(type: 'string', length: 255)]
    private string $name;

    #[Assert\GreaterThanOrEqual(value: 0, message: 'Position should be greater than or equal to 0.')]
    #[Assert\NotBlank(message: 'Position is required.')]
    #[ORM\Column(type: 'integer')]
    private int $position;

    #[Assert\NotBlank(message: 'Type is required.')]
    #[ORM\Column(type: 'string', length: 7, enumType: TransactionType::class)]
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
