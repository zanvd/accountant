<?php

namespace App\Entity;

use App\Repository\CategoryRepository;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\Mapping as ORM;

#[ORM\Entity(repositoryClass: CategoryRepository::class)]
#[ORM\Table(name: 'category')]
class Category
{
    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column(type: 'integer')]
    private int $id;

    #[ORM\Column(length: 7, type: 'string')]
    private string $color = '';

    #[ORM\Column(length: 150, type: 'string')]
    private string $description = '';

    #[ORM\Column(length: 30, type: 'string')]
    private string $name;

    #[ORM\Column(length: 7, type: 'string')]
    private string $text_color = '#000000';

    #[ORM\OneToMany(mappedBy: 'category', targetEntity: Transaction::class)]
    private Collection $transactions;

    #[ORM\OneToMany(mappedBy: 'category', targetEntity: TransactionTemplate::class)]
    private $transactionTemplates;

    #[ORM\ManyToOne(targetEntity: User::class, inversedBy: 'categories')]
    #[ORM\JoinColumn(nullable: false)]
    private $user;

    public function __construct()
    {
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

    public function getTextColor(): ?string
    {
        return $this->text_color;
    }

    public function setTextColor(string $text_color): self
    {
        $this->text_color = $text_color;

        return $this;
    }

    /**
     * @return Collection<int, Transaction>
     */
    public function getTransactions(): Collection
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
     * @return Collection<int, TransactionTemplate>
     */
    public function getTransactionTemplates(): Collection
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
