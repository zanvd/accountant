<?php

namespace App\Entity;

use App\Enum\UserValidationType;
use App\Repository\UserRepository;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\Mapping as ORM;
use Symfony\Bridge\Doctrine\Validator\Constraints\UniqueEntity;
use Symfony\Component\Security\Core\User\PasswordAuthenticatedUserInterface;
use Symfony\Component\Security\Core\User\UserInterface;
use Symfony\Component\Validator\Constraints as Assert;

#[ORM\Entity(repositoryClass: UserRepository::class)]
#[UniqueEntity(fields: ['username'], message: 'There is already an account with this username')]
class User implements UserInterface, PasswordAuthenticatedUserInterface
{
    #[ORM\Id]
    #[ORM\GeneratedValue]
    #[ORM\Column(type: 'integer')]
    private int $id;

    #[Assert\NotBlank(message: 'Email is required.', groups: [UserValidationType::Register->value])]
    #[ORM\Column(type: 'string', length: 255)]
    private string $email;

    #[ORM\Column(type: 'boolean')]
    private bool $isVerified = false;

    #[Assert\Length(
        min: 6,
        max: 4096, // Max length allowed by Symfony for security reasons.
        minMessage: 'Your password should have at least {{ limit }} characters.',
        maxMessage: 'Password cannot have more than {{ limit }} characters.',
        groups: [UserValidationType::Register->value, UserValidationType::Reset->value]
    )]
    #[Assert\NotBlank(
        message: 'Password is required.',
        groups: [
            UserValidationType::Login->value,
            UserValidationType::Register->value,
            UserValidationType::Reset->value,
        ],
    )]
    #[ORM\Column(type: 'string')]
    private string $password;

    #[ORM\Column(type: 'json')]
    private array $roles = [];

    #[Assert\NotBlank(
        message: 'Username is required.',
        groups: [
            UserValidationType::Login->value,
            UserValidationType::ForgotPassword->value,
            UserValidationType::Register->value,
        ],
    )]
    #[ORM\Column(type: 'string', length: 255, unique: true)]
    private string $username;

    #[ORM\OneToMany(mappedBy: 'user', targetEntity: Category::class, orphanRemoval: true)]
    private Collection $categories;

    #[ORM\OneToMany(mappedBy: 'user', targetEntity: RecurringTransaction::class, orphanRemoval: true)]
    private Collection $recurringTransactions;

    #[ORM\OneToMany(mappedBy: 'user', targetEntity: Transaction::class, orphanRemoval: true)]
    private Collection $transactions;

    #[ORM\OneToMany(mappedBy: 'user', targetEntity: TransactionTemplate::class, orphanRemoval: true)]
    private Collection $transactionTemplates;

    public function __construct()
    {
        $this->categories = new ArrayCollection();
        $this->recurringTransactions = new ArrayCollection();
        $this->transactions = new ArrayCollection();
        $this->transactionTemplates = new ArrayCollection();
    }

    public function getId(): ?int
    {
        return $this->id;
    }

    public function getEmail(): string
    {
        return $this->email;
    }

    public function setEmail(string $email): self
    {
        $this->email = $email;

        return $this;
    }

    public function isVerified(): bool
    {
        return $this->isVerified;
    }

    public function setIsVerified(bool $isVerified): self
    {
        $this->isVerified = $isVerified;

        return $this;
    }

    /**
     * @see UserInterface
     */
    public function eraseCredentials()
    {
        // If you store any temporary, sensitive data on the user, clear it here
        // $this->plainPassword = null;
    }

    /**
     * @see PasswordAuthenticatedUserInterface
     */
    public function getPassword(): string
    {
        return $this->password;
    }

    public function setPassword(string $password): self
    {
        $this->password = $password;

        return $this;
    }

    /**
     * A visual identifier that represents this user.
     *
     * @see UserInterface
     */
    public function getUserIdentifier(): string
    {
        return $this->username;
    }

    /**
     * @see UserInterface
     */
    public function getRoles(): array
    {
        $roles = $this->roles;
        // guarantee every user at least has ROLE_USER
        $roles[] = 'ROLE_USER';

        return array_unique($roles);
    }

    public function setRoles(array $roles): self
    {
        $this->roles = $roles;

        return $this;
    }

    public function getUsername(): string
    {
        return $this->username;
    }

    public function setUsername(string $username): self
    {
        $this->username = $username;

        return $this;
    }

    /**
     * @return Collection<int, Category>
     */
    public function getCategories(): Collection
    {
        return $this->categories;
    }

    public function addCategory(Category $category): self
    {
        if (!$this->categories->contains($category)) {
            $this->categories[] = $category;
            $category->setUser($this);
        }

        return $this;
    }

    public function removeCategory(Category $category): self
    {
        if ($this->categories->removeElement($category)) {
            // set the owning side to null (unless already changed)
            if ($category->getUser() === $this) {
                $category->setUser(null);
            }
        }

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
            $recurringTransaction->setUser($this);
        }

        return $this;
    }

    public function removeRecurringTransaction(RecurringTransaction $recurringTransaction): self
    {
        if ($this->recurringTransactions->removeElement($recurringTransaction)) {
            // set the owning side to null (unless already changed)
            if ($recurringTransaction->getUser() === $this) {
                $recurringTransaction->setUser(null);
            }
        }

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
            $transaction->setUser($this);
        }

        return $this;
    }

    public function removeTransaction(Transaction $transaction): self
    {
        if ($this->transactions->removeElement($transaction)) {
            // set the owning side to null (unless already changed)
            if ($transaction->getUser() === $this) {
                $transaction->setUser(null);
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
            $transactionTemplate->setUser($this);
        }

        return $this;
    }

    public function removeTransactionTemplate(TransactionTemplate $transactionTemplate): self
    {
        if ($this->transactionTemplates->removeElement($transactionTemplate)) {
            // set the owning side to null (unless already changed)
            if ($transactionTemplate->getUser() === $this) {
                $transactionTemplate->setUser(null);
            }
        }

        return $this;
    }
}
