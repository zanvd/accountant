<?php

namespace App\Exporter;

use App\Entity\Transaction;
use App\Repository\TransactionRepository;
use Doctrine\Persistence\ManagerRegistry;
use Doctrine\Persistence\ObjectRepository;
use Symfony\Component\Serializer\Encoder\CsvEncoder;

class Csv implements ExporterInterface
{
    private ObjectRepository|TransactionRepository $trRepo;

    public function __construct(ManagerRegistry $doctrine)
    {
        $this->trRepo = $doctrine->getRepository(Transaction::class);
    }

    public function export(): string
    {
        $transactions = $this->trRepo->findAll();

        $encoder = new CsvEncoder();

        return $encoder->encode(
            array_map(fn(Transaction $t) => [
                'name' => $t->getName(),
                'category' => $t->getCategory()->getName(),
                'amount' => $t->getAmount(),
                'date' => $t->getTransactionDate()->format('d/m/Y'),
                'summary' => $t->getSummary(),
            ], $transactions),
            'csv'
        );
    }
}
