<?php

namespace App\Controller\Private;

use App\Entity\RecurringTransaction;
use App\Form\Type\RecurringTransactionType;
use App\Messenger\RecurringTransactionSender;
use App\Repository\RecurringTransactionRepository;
use DateTime;
use Doctrine\Persistence\ManagerRegistry;
use Doctrine\Persistence\ObjectRepository;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/recurring-transaction', name: 'recurring_transaction_')]
class RecurringTransactionController extends AbstractController
{
    private ManagerRegistry $doctrine;
    private ObjectRepository|RecurringTransactionRepository $recTrRepo;
    private RecurringTransactionSender $msgSender;

    public function __construct(ManagerRegistry $doctrine, RecurringTransactionSender $msgSender)
    {
        $this->doctrine = $doctrine;
        $this->recTrRepo = $doctrine->getRepository(RecurringTransaction::class);
        $this->msgSender = $msgSender;
    }

    #[Route('/add', name: 'add', methods: ['GET', 'POST'])]
    public function add(Request $request): Response
    {
        $rt = new RecurringTransaction();
        $rt->setStartDate(new DateTime())
            ->setUser($this->getUser());

        $f = $this->createForm(RecurringTransactionType::class, $rt);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $this->recTrRepo->add($rt, true);

            if (!$this->msgSender->sendAsyncMessage($rt)) {
                $this->addFlash('rec_trans_error', 'Failed to set up a recurring transaction. Please try again.');
            } else {
                return $this->redirectToRoute('recurring_transaction_index');
            }
        }

        return $this->renderForm('recurringTransaction/add.html.twig', ['form' => $f]);
    }

    #[Route('/delete/{id}', name: 'delete', methods: ['GET'])]
    public function delete(int $id): Response
    {
        $rt = $this->recTrRepo->findOneBy(['id' => $id, 'user' => $this->getUser()]);
        if (!$rt) {
            throw $this->createNotFoundException("No recurring transaction with id $id.");
        }

        $this->recTrRepo->remove($rt, true);

        return $this->redirectToRoute('recurring_transaction_index');
    }

    #[Route('/edit/{id}', name: 'edit', methods: ['GET', 'POST'])]
    public function edit(int $id, Request $request): Response
    {
        $rt = $this->recTrRepo->findOneBy(['id' => $id, 'user' => $this->getUser()]);
        if (!$rt) {
            throw $this->createNotFoundException("No recurring transaction with id $id.");
        }

        $f = $this->createForm(RecurringTransactionType::class, $rt);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $this->doctrine->getManager()->flush();

            if (!$this->msgSender->sendAsyncMessage($rt)) {
                $this->addFlash('rec_trans_error', 'Failed to set up a recurring transaction. Please try again.');
            } else {
                return $this->redirectToRoute('recurring_transaction_index');
            }
        }

        return $this->renderForm('recurringTransaction/edit.html.twig', [
            'form' => $f,
            'recurringTransaction' => $rt,
        ]);
    }

    #[Route('', name: 'index', methods: ['GET'])]
    public function index(): Response
    {
        return $this->render('recurringTransaction/index.html.twig', [
            'recurringTransactions' => $this->recTrRepo->findBy(['user' => $this->getUser()]),
        ]);
    }

    #[Route('/view/{id}', name: 'view', methods: ['GET'])]
    public function view(int $id): Response
    {
        $rt = $this->recTrRepo->findOneBy(['id' => $id, 'user' => $this->getUser()]);
        if (!$rt) {
            throw $this->createNotFoundException("No recurring transaction with id $id.");
        }

        return $this->render('recurringTransaction/view.html.twig', ['recurringTransaction' => $rt]);
    }
}
