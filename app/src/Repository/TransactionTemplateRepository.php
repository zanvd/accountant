<?php

namespace App\Repository;

use App\Entity\TransactionTemplate;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;

/**
 * @extends ServiceEntityRepository<TransactionTemplate>
 *
 * @method TransactionTemplate|null find($id, $lockMode = null, $lockVersion = null)
 * @method TransactionTemplate|null findOneBy(array $criteria, array $orderBy = null)
 * @method TransactionTemplate[]    findAll()
 * @method TransactionTemplate[]    findBy(array $criteria, array $orderBy = null, $limit = null, $offset = null)
 */
class TransactionTemplateRepository extends ServiceEntityRepository
{
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, TransactionTemplate::class);
    }

    public function add(TransactionTemplate $entity, bool $flush = false): void
    {
        $this->getEntityManager()->persist($entity);

        if ($flush) {
            $this->getEntityManager()->flush();
        }
    }

    public function remove(TransactionTemplate $entity, bool $flush = false): void
    {
        $this->getEntityManager()->remove($entity);

        if ($flush) {
            $this->getEntityManager()->flush();
        }
    }

//    /**
//     * @return TransactionTemplate[] Returns an array of TransactionTemplate objects
//     */
//    public function findByExampleField($value): array
//    {
//        return $this->createQueryBuilder('t')
//            ->andWhere('t.exampleField = :val')
//            ->setParameter('val', $value)
//            ->orderBy('t.id', 'ASC')
//            ->setMaxResults(10)
//            ->getQuery()
//            ->getResult()
//        ;
//    }

//    public function findOneBySomeField($value): ?TransactionTemplate
//    {
//        return $this->createQueryBuilder('t')
//            ->andWhere('t.exampleField = :val')
//            ->setParameter('val', $value)
//            ->getQuery()
//            ->getOneOrNullResult()
//        ;
//    }
}
