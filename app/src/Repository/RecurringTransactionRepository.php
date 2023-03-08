<?php

namespace App\Repository;

use App\Entity\RecurringTransaction;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;

/**
 * @extends ServiceEntityRepository<RecurringTransaction>
 *
 * @method RecurringTransaction|null find($id, $lockMode = null, $lockVersion = null)
 * @method RecurringTransaction|null findOneBy(array $criteria, array $orderBy = null)
 * @method RecurringTransaction[]    findAll()
 * @method RecurringTransaction[]    findBy(array $criteria, array $orderBy = null, $limit = null, $offset = null)
 */
class RecurringTransactionRepository extends ServiceEntityRepository
{
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, RecurringTransaction::class);
    }

    public function add(RecurringTransaction $entity, bool $flush = false): void
    {
        $this->getEntityManager()->persist($entity);
        if ($flush) {
            $this->getEntityManager()->flush();
        }
    }

    public function remove(RecurringTransaction $entity, bool $flush = false): void
    {
        $this->getEntityManager()->remove($entity);
        if ($flush) {
            $this->getEntityManager()->flush();
        }
    }
}
