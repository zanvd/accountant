<?php

namespace App\Repository;

use App\Entity\Transaction;
use DateTime;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Common\Collections\Criteria;
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
    public function __construct(ManagerRegistry $registry)
    {
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
    public function listForPeriod(UserInterface $user, ?DateTime $from = null, ?DateTime $to = null): array
    {
        $from ??= new DateTime('first day of this month');
        $to ??= new DateTime('last day of this month');

        $criteria = new Criteria();
        $criteria
            ->where(
                Criteria::expr()->andX(
                    Criteria::expr()->gte('transactionDate', $from),
                    Criteria::expr()->lte('transactionDate', $to),
                    Criteria::expr()->eq('user', $user)
                )
            )
            ->orderBy([
                'transactionDate' => 'DESC',
                'id' => 'DESC',
            ]);

        return $this->matching($criteria)->toArray();
    }

    public function remove(Transaction $entity, bool $flush = false): void
    {
        $this->getEntityManager()->remove($entity);

        if ($flush) {
            $this->getEntityManager()->flush();
        }
    }
}
