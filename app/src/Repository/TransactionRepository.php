<?php

namespace App\Repository;

use App\Entity\Category;
use App\Entity\Transaction;
use App\Filter\DateTimeFilter;
use App\Filter\EntityFilter;
use App\Filter\TransactionFilter;
use DateTime;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;
use Symfony\Component\Security\Core\User\UserInterface;

/**
 * @extends ServiceEntityRepository<Transaction>
 *
 * @method Transaction|null find($id, $lockMode = null, $lockVersion = null)
 * @method Transaction|null findOneBy(array $criteria, array $orderBy = null)
 * @method Transaction[]    findAll()
 * @method Transaction[]    findBy(array $criteria, array $orderBy = null, $limit = null, $offset = null)
 */
class TransactionRepository extends ServiceEntityRepository
{
    private DateTimeFilter $dateTimeFilter;
    private EntityFilter $entityFilter;

    public function __construct(ManagerRegistry $registry, DateTimeFilter $dateTimeFilter, EntityFilter $entityFilter)
    {
        $this->dateTimeFilter = $dateTimeFilter;
        $this->entityFilter = $entityFilter;

        parent::__construct($registry, Transaction::class);
    }

    public function add(Transaction $entity, bool $flush = false): void
    {
        $this->getEntityManager()->persist($entity);

        if ($flush) {
            $this->getEntityManager()->flush();
        }
    }

    /**
     * @return Transaction[]
     */
    public function getFilteredTransactions(
        UserInterface $user,
        ?Category $category = null,
        ?DateTime $from = null,
        ?DateTime $to = null
    ): array {
        $from ??= new DateTime('first day of this month');
        $from->setTime(0, 0);
        $to ??= new DateTime('last day of this month');
        $to->setTime(23, 59, 59);

        $qb = $this->createQueryBuilder('t');
        $qb->andWhere($qb->expr()->eq('t.user', ':user'))
            ->setParameter('user', $user)
            ->orderBy('t.transactionDate', 'DESC')
            ->addOrderBy('t.id', 'DESC');
        $qb = $this->dateTimeFilter->apply($qb, 't', 'transactionDate', 'transactionDate', $from, $to);
        if ($category) {
            $qb = $this->entityFilter->apply($qb, 't', 'category', $category);
        }

        return $qb->getQuery()->getResult();
    }

    public function remove(Transaction $entity, bool $flush = false): void
    {
        $this->getEntityManager()->remove($entity);

        if ($flush) {
            $this->getEntityManager()->flush();
        }
    }
}
