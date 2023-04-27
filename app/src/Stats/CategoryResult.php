<?php

namespace App\Stats;

class CategoryResult
{
    private int $id;

    private int $amount;

    private string $color;

    private string $name;

    public function __construct(int $id, int $amount, string $color, string $name)
    {
        $this->id = $id;
        $this->amount = $amount;
        $this->color = $color;
        $this->name = $name;
    }

    public function getId(): int
    {
        return $this->id;
    }

    public function getAmount(): int
    {
        return $this->amount;
    }

    public function getColor(): string
    {
        return $this->color;
    }

    public function getName(): string
    {
        return $this->name;
    }
}