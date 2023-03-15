<?php

namespace App\Filter;

use Doctrine\ORM\QueryBuilder;

class EntityFilter
{
    public function apply(QueryBuilder $qb, string $alias, string $field, object $entity): QueryBuilder
    {
        return $qb
            ->andWhere($qb->expr()->eq("$alias.$field", ":$field"))
            ->setParameter(":$field", $entity);
    }
}
