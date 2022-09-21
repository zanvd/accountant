<?php

namespace App\Balance;

use App\Entity\Transaction;
use App\Repository\TransactionRepository;
use DateTime;
use Doctrine\Persistence\ManagerRegistry;
use Doctrine\Persistence\ObjectRepository;
use Symfony\Component\Security\Core\User\UserInterface;

class Calculator
{
    private ObjectRepository|TransactionRepository $trRepo;

    public function __construct(ManagerRegistry $doctrine)
    {
        $this->trRepo = $doctrine->getRepository(Transaction::class);
    }

    public function getAllTimeBalance(UserInterface $user): Balance
    {
        $b = new Balance();
        foreach ($this->trRepo->findBy(['user' => $user]) as $transaction) {
            if (($amount = $transaction->getAmount()) > 0) {
                $b->setIncome($b->getIncome() + $amount);
            } else {
                $b->setOutcome($b->getOutcome() + $amount);
            }
        }
        $b->setSavings($b->getIncome() + $b->getOutcome());

        return $b;
    }

    public function getPeriodBalance(UserInterface $user, ?DateTime $from = null, ?DateTime $to = null): Balance
    {
        $b = new Balance();
        foreach ($this->trRepo->getTransactionsForPeriod($user, $from, $to) as $transaction) {
            if (($amount = $transaction->getAmount()) > 0) {
                $b->setIncome($b->getIncome() + $amount);
            } else {
                $b->setOutcome($b->getOutcome() + $amount);
            }
        }
        $b->setSavings($b->getIncome() + $b->getOutcome());

        return $b;
    }
}
