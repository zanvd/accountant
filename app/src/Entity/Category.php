<?php

namespace App\Entity;

use App\Repository\CategoryRepository;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\Mapping as ORM;
use Symfony\Component\Cache\Adapter\ArrayAdapter;

#[ORM\Entity(repositoryClass: CategoryRepository::class)]
#[ORM\Table(name: 'category')]
class Category
{
    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column(type: 'integer')]
    private int $id;

    #[ORM\Column(type: 'string', length: 7)]
    private string $color = '';

    #[ORM\Column(type: 'string', length: 150)]
    private string $description = '';

    #[ORM\Column(type: 'string', length: 30)]
    private string $name;

    #[ORM\OneToMany(mappedBy: 'category', targetEntity: RecurringTransaction::class)]
    private Collection $recurringTransactions;

    #[ORM\Column(type: 'string', length: 7)]
    private string $textColor = '#000000';

    #[ORM\OneToMany(mappedBy: 'category', targetEntity: Transaction::class)]
    private Collection $transactions;

    #[ORM\OneToMany(mappedBy: 'category', targetEntity: TransactionTemplate::class)]
    private Collection $transactionTemplates;

    #[ORM\ManyToOne(targetEntity: User::class, inversedBy: 'categories')]
    #[ORM\JoinColumn(nullable: false)]
    private ?User $user;

    public function __construct()
    {
        $this->recurringTransactions = new ArrayCollection();
        $this->transactions = new ArrayCollection();
        $this->transactionTemplates = new ArrayCollection();
    }

    public function getId(): ?int
    {
        return $this->id;
    }

    public function getColor(): ?string
    {
        return $this->color;
    }

    public function setColor(string $color): self
    {
        $this->color = $color;

        return $this;
    }

    public function getDescription(): ?string
    {
        return $this->description;
    }

    public function setDescription(string $description): self
    {
        $this->description = $description;

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

    /**
     * @return ArrayCollection<int, RecurringTransaction>
     */
    public function getRecurringTransactions(): ArrayCollection
    {
        return $this->transactions;
    }

    public function addRecurringTransaction(RecurringTransaction $recurringTransaction): self
    {
        if (!$this->recurringTransactions->contains($recurringTransaction)) {
            $this->recurringTransactions[] = $recurringTransaction;
            $recurringTransaction->setCategory($this);
        }

        return $this;
    }

    public function removeRecurringTransaction(RecurringTransaction $recurringTransaction): self
    {
        if ($this->recurringTransactions->removeElement($recurringTransaction)) {
            // set the owning side to null (unless already changed)
            if ($recurringTransaction->getCategory() === $this) {
                $recurringTransaction->setCategory(null);
            }
        }

        return $this;
    }

    public function getTextColor(): ?string
    {
        return $this->textColor;
    }

    public function setTextColor(string $textColor): self
    {
        $this->textColor = $textColor;

        return $this;
    }

    /**
     * @return ArrayCollection<int, Transaction>
     */
    public function getTransactions(): ArrayCollection
    {
        return $this->transactions;
    }

    public function addTransaction(Transaction $transaction): self
    {
        if (!$this->transactions->contains($transaction)) {
            $this->transactions[] = $transaction;
            $transaction->setCategory($this);
        }

        return $this;
    }

    public function removeTransaction(Transaction $transaction): self
    {
        if ($this->transactions->removeElement($transaction)) {
            // set the owning side to null (unless already changed)
            if ($transaction->getCategory() === $this) {
                $transaction->setCategory(null);
            }
        }

        return $this;
    }

    /**
     * @return ArrayCollection<int, TransactionTemplate>
     */
    public function getTransactionTemplates(): ArrayCollection
    {
        return $this->transactionTemplates;
    }

    public function addTransactionTemplate(TransactionTemplate $transactionTemplate): self
    {
        if (!$this->transactionTemplates->contains($transactionTemplate)) {
            $this->transactionTemplates[] = $transactionTemplate;
            $transactionTemplate->setCategory($this);
        }

        return $this;
    }

    public function removeTransactionTemplate(TransactionTemplate $transactionTemplate): self
    {
        if ($this->transactionTemplates->removeElement($transactionTemplate)) {
            // set the owning side to null (unless already changed)
            if ($transactionTemplate->getCategory() === $this) {
                $transactionTemplate->setCategory(null);
            }
        }

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
}
