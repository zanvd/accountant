<?php

namespace App\Stats;

use App\Entity\Category as CategoryEntity;
use App\Repository\CategoryRepository;
use Doctrine\Persistence\ObjectRepository;
use Doctrine\Persistence\ManagerRegistry;

class Category
{
    private ObjectRepository|CategoryRepository $catRepo;

    public function __construct(ManagerRegistry $doctrine)
    {
        $this->catRepo = $doctrine->getRepository(CategoryEntity::class);
    }

    /**
     * @return array []CategoryResult
     */
    public function calcAll(): array
    {
        $qb = $this->catRepo->createQueryBuilder('c');
        $qb->select('c.id', 'c.color', 'c.name', 'SUM(t.amount) AS amount')
            ->innerJoin('c.transactions', 't', 'WITH', 't.category = c')
            ->groupBy('c');

        $res = [];
        foreach ($qb->getQuery()->getResult() as $r) {
            $res[] = new CategoryResult($r['id'], $r['amount'], $r['color'], $r['name']);
        }

        return $res;
    }
}
