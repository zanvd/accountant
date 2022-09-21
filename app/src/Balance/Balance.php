<?php

namespace App\Balance;

class Balance
{
    private ?string $income = null;
    private ?string $outcome = null;
    private ?string $savings = null;

    public function getIncome(): ?string
    {
        return $this->income;
    }

    public function setIncome(?string $income): self
    {
        $this->income = $income;

        return $this;
    }

    public function getOutcome(): ?string
    {
        return $this->outcome;
    }

    public function setOutcome(?string $outcome): self
    {
        $this->outcome = $outcome;

        return $this;
    }

    public function getSavings(): ?string
    {
        return $this->savings;
    }

    public function setSavings(?string $savings): self
    {
        $this->savings = $savings;

        return $this;
    }
}
