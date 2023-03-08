<?php

namespace App\Messenger;

use App\Entity\RecurringTransaction;
use DateTime;
use DateTimeInterface;
use Psr\Log\LoggerInterface;
use Symfony\Component\Messenger\Envelope;
use Symfony\Component\Messenger\MessageBusInterface;
use Symfony\Component\Messenger\Stamp\DelayStamp;

class RecurringTransactionSender
{
    private MessageBusInterface $bus;
    private LoggerInterface $logger;

    public function __construct(MessageBusInterface $bus, LoggerInterface $logger)
    {
        $this->bus = $bus;
        $this->logger = $logger;
    }

    public function sendAsyncMessage(RecurringTransaction $recTrans): bool
    {
        if (!is_null($nextDate = $this->getNextOccurrence($recTrans))) {
            $this->bus->dispatch(new Envelope(
                new RecurringTransactionMessage($recTrans->getId()),
                [
                    DelayStamp::delayUntil($nextDate)
                ]
            ));

            return true;
        }

        $this->logger->error("Failed to get the next occurrence for {$recTrans->getId()}.");

        return false;
    }

    private function getNextOccurrence(RecurringTransaction $recTrans): ?DateTimeInterface
    {
        $today = new DateTime();
        if ($recTrans->getStartDate() > $today) {
            $next = $recTrans->getStartDate();
        } else {
            $next = $today->modify("+{$recTrans->getPeriodNum()} {$recTrans->getPeriodType()->value}");
        }

        return is_null($recTrans->getEndDate()) || $next <= $recTrans->getEndDate() ? $next : null;
    }
}
