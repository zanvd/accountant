<?php

namespace App\Messenger;

class RecurringTransactionMessage implements MessageInterface
{
    private string $content;

    public function __construct(int $recTransId)
    {
        $this->content = (string)$recTransId;
    }

    public function getContent(): int
    {
        return (int)$this->content;
    }
}
