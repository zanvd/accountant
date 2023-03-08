<?php

namespace App\Messenger;

use App\Entity\RecurringTransaction;
use App\Entity\Transaction;
use App\Repository\RecurringTransactionRepository;
use App\Repository\TransactionRepository;
use DateInterval;
use DatePeriod;
use DateTime;
use DateTimeInterface;
use Doctrine\ORM\EntityNotFoundException;
use Doctrine\Persistence\ManagerRegistry;
use Doctrine\Persistence\ObjectRepository;
use Exception;
use Psr\Log\LoggerInterface;
use Symfony\Component\Messenger\Attribute\AsMessageHandler;
use Symfony\Component\Messenger\Exception\UnrecoverableMessageHandlingException;

#[AsMessageHandler]
class RecurringTransactionHandler
{
    private LoggerInterface $logger;
    private RecurringTransactionSender $msgSender;
    private ObjectRepository|RecurringTransactionRepository $recTrRepo;
    private ObjectRepository|TransactionRepository $trRepo;

    public function __construct(
        ManagerRegistry $doctrine,
        LoggerInterface $logger,
        RecurringTransactionSender $msgSender
    ) {
        $this->msgSender = $msgSender;
        $this->logger = $logger;
        $this->recTrRepo = $doctrine->getRepository(RecurringTransaction::class);
        $this->trRepo = $doctrine->getRepository(RecurringTransaction::class);
    }

    /**
     * @throws Exception
     */
    public function __invoke(RecurringTransactionMessage $message): void
    {
        $recTrans = $this->recTrRepo->find($message->getContent());
        if (!$recTrans) {
            $this->log("Failed to retrieve transaction {$message->getContent()}.");

            throw EntityNotFoundException::fromClassNameAndIdentifier(
                $this->recTrRepo->getClassName(),
                [(string)$message->getContent()]
            );
        }

        if (!$this->shouldOccurToday($recTrans)) {
            $this->log("Transactions should not occur today: {$recTrans->getId()}.");

            throw new UnrecoverableMessageHandlingException("Transactions should not occur today: {$recTrans->getId()}.");
        }

        $trans = (new Transaction())
            ->setAmount($recTrans->getAmount())
            ->setCategory($recTrans->getCategory())
            ->setName($recTrans->getName())
            ->setSummary($recTrans->getSummary())
            ->setTransactionDate(new DateTime())
            ->setUser($recTrans->getUser());
        $this->trRepo->add($trans, true);

        $this->msgSender->sendAsyncMessage($recTrans);
    }

    private function shouldOccurToday(RecurringTransaction $recTrans): bool
    {
        try {
            $endDate = $this->getEndDate($recTrans);
        } catch (Exception $e) {
            $this->log("Failed to obtain end date: {$e->getMessage()}");
            return false;
        }
        $today = new Datetime();
        if ($today < $recTrans->getStartDate() || $today > $endDate) {
            return false;
        }

        try {
            $period = new DatePeriod(
                $recTrans->getStartDate(),
                new DateInterval("P{$recTrans->getPeriodNum()}{$recTrans->getPeriodType()->getIntervalSymbol()}"),
                $endDate,
                DatePeriod::EXCLUDE_START_DATE | DatePeriod::INCLUDE_END_DATE
            );
        } catch (Exception $e) {
            $this->log("Failed to create a period: {$e->getMessage()}");
            return false;
        }

        return in_array($today, iterator_to_array($period));
    }

    /**
     * Returns end date if not null, otherwise today + set period.
     * @throws Exception
     */
    private function getEndDate(RecurringTransaction $recTrans): DateTimeInterface
    {
        return !is_null($recTrans->getEndDate())
            ? $recTrans->getEndDate()
            : (new DateTime())
                ->modify("+{$recTrans->getPeriodNum()} {$recTrans->getPeriodType()->value}");
    }

    private function log(string $msg, string $type='error'): void {
        $msg = '[' . self::class . "] $msg";
        match ($type) {
            'debug' => $this->logger->debug($msg),
            'error' => $this->logger->error($msg),
            'info' => $this->logger->info($msg),
        };
    }
}
