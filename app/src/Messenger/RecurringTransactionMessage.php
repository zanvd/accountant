<?php

namespace App\Messenger;

use DateTimeInterface;
use Symfony\Component\Messenger\Exception\UnrecoverableMessageHandlingException;

class RecurringTransactionMessage implements MessageInterface
{
    private string $content;

    public function __construct(DateTimeInterface $modifiedAt, int $recTransId)
    {
        $this->content = serialize([
            'id' => $recTransId,
            'modified_at' => $modifiedAt,
        ]);
    }

    public function getContent(): int
    {
        $content = unserialize($this->content);
        if ($content === false) {
            throw new UnrecoverableMessageHandlingException("Failed to unserialize data: $this->content");
        }
        return $content['id'];
    }

    public function getModifiedAt(): DateTimeInterface
    {
        $content = unserialize($this->content);
        if ($content === false) {
            throw new UnrecoverableMessageHandlingException("Failed to unserialize data: $this->content");
        }
        return $content['modified_at'];
    }
}
