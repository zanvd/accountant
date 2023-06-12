<?php

namespace App\Entity;

use App\Enum\RecurringPeriodType;
use App\Repository\RecurringTransactionRepository;
use ContainerI02yXV4\getTwig_Runtime_HttpkernelService;
use DateTime;
use DateTimeInterface;
use Doctrine\ORM\Mapping as ORM;
use Symfony\Component\Validator\Constraints as Assert;

#[ORM\Entity(repositoryClass: RecurringTransactionRepository::class)]
#[ORM\HasLifecycleCallbacks]
class RecurringTransaction
{
    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column(type: 'integer')]
    private int $id;

    #[Assert\NotBlank(message: 'Amount is required.')]
    #[ORM\Column(type: 'decimal', precision: 19, scale: 4)]
    private string $amount;

    #[ORM\ManyToOne(targetEntity: Category::class, inversedBy: 'recurringTransactions')]
    #[ORM\JoinColumn(nullable: false)]
    private ?Category $category;

    #[ORM\Column(type: 'date', nullable: true)]
    private ?DateTimeInterface $endDate = null;

    #[ORM\Column(type: 'datetime')]
    private ?DateTimeInterface $modifiedAt = null;

    #[Assert\NotBlank(message: 'Name is required.')]
    #[ORM\Column(type: 'string', length: 255)]
    private string $name;

    #[Assert\GreaterThan(value: 0, message: 'Period should be greater than 0.')]
    #[Assert\NotBlank(message: 'Period is required.')]
    #[ORM\Column(type: 'integer', nullable: false)]
    private int $periodNum;

    #[ORM\Column(type: 'string', length: 5, enumType: RecurringPeriodType::class)]
    private RecurringPeriodType $periodType;

    #[Assert\GreaterThanOrEqual('today')]
    #[Assert\NotBlank(message: 'Start date is required.')]
    #[Assert\LessThan(propertyPath: 'endDate')]
    #[ORM\Column(type: 'date', nullable: false)]
    private DateTimeInterface $startDate;

    #[Assert\NotNull]
    #[ORM\Column(type: 'string', length: 255)]
    private string $summary = '';

    #[ORM\ManyToOne(targetEntity: User::class, inversedBy: 'recurringTransactions')]
    #[ORM\JoinColumn(nullable: false)]
    private ?User $user;

    public function getId(): int
    {
        return $this->id;
    }

    public function getAmount(): string
    {
        return $this->amount;
    }

    public function setAmount(string $amount): void
    {
        $this->amount = $amount;
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

    public function getEndDate(): ?DateTimeInterface
    {
        return $this->endDate;
    }

    public function setEndDate(?DateTimeInterface $endDate = null): self
    {
        $this->endDate = $endDate;
        return $this;
    }

    public function getModifiedAt(): ?DateTimeInterface
    {
        return $this->modifiedAt;
    }

    public function setModifiedAt(?DateTimeInterface $modifiedAt = null): self
    {
        $this->modifiedAt = $modifiedAt;
        return $this;
    }

    public function getName(): string
    {
        return $this->name;
    }

    public function setName(string $name): self
    {
        $this->name = $name;
        return $this;
    }

    public function getPeriodNum(): int
    {
        return $this->periodNum;
    }

    public function setPeriodNum(int $periodNum): self
    {
        $this->periodNum = $periodNum;
        return $this;
    }

    public function getPeriodType(): RecurringPeriodType
    {
        return $this->periodType;
    }

    public function setPeriodType(RecurringPeriodType $periodType): self
    {
        $this->periodType = $periodType;
        return $this;
    }

    public function getStartDate(): DateTimeInterface
    {
        return $this->startDate;
    }

    public function setStartDate(DateTimeInterface $startDate): self
    {
        $this->startDate = $startDate;
        return $this;
    }

    public function getSummary(): string
    {
        return $this->summary;
    }

    public function setSummary(string $summary): self
    {
        $this->summary = $summary;
        return $this;
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

    #[ORM\PrePersist]
    #[ORM\PreUpdate]
    public function onUpdate(): void
    {
        $this->setModifiedAt(new DateTime('now'));
    }
}
