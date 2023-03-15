<?php

namespace App\Filter;

use DateTime;
use Doctrine\ORM\QueryBuilder;

class DateTimeFilter
{
    public function apply(
        QueryBuilder $qb,
        string $alias,
        string $fromField,
        string $toField,
        DateTime $from,
        DateTime $to
    ): QueryBuilder {
        return $qb
            ->andWhere(
                $qb->expr()->andX(
                    $qb->expr()->gte("$alias.$fromField", ":from"),
                    $qb->expr()->lte("$alias.$toField", ":to"),
                ),
            )
            ->setParameter('from', $from)
            ->setParameter('to', $to);
    }
}
