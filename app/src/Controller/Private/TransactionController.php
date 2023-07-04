<?php

namespace App\Controller\Private;

use App\Balance\Calculator;
use App\Entity\Category;
use App\Entity\Transaction;
use App\Enum\TransactionType as TransactionTypeEnum;
use App\Form\Type\TransactionType;
use App\Repository\CategoryRepository;
use App\Repository\TransactionRepository;
use DateTime;
use Doctrine\Persistence\ManagerRegistry;
use Doctrine\Persistence\ObjectRepository;
use Exception;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/transaction', name: 'transaction_')]
class TransactionController extends AbstractController
{
    private Calculator $calculator;
    private ObjectRepository|CategoryRepository $catRepo;
    private ManagerRegistry $doctrine;
    private ObjectRepository|TransactionRepository $trRepo;

    public function __construct(Calculator $calculator, ManagerRegistry $doctrine)
    {
        $this->calculator = $calculator;
        $this->doctrine = $doctrine;
        $this->catRepo = $this->doctrine->getRepository(Category::class);
        $this->trRepo = $this->doctrine->getRepository(Transaction::class);
    }

    #[Route('/add', name: 'add', methods: ['GET', 'POST'])]
    public function add(Request $request): Response
    {
        $c = null;
        if ($cid = $request->query->get('category')) {
            $c = $this->catRepo->findOneBy(['id' => $cid, 'user' => $this->getUser()]);
        }

        $t = new Transaction();
        $t->setCategory($c)
            ->setName($request->query->get('name', ''))
            ->setTransactionDate(new DateTime())
            ->setUser($this->getUser());
        if (TransactionTypeEnum::tryFrom($request->query->get('type')) === TransactionTypeEnum::Outcome) {
            // TODO: Consider alternative way.
            $t->setAmount('-1');
        }

        $f = $this->createForm(TransactionType::class, $t);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $this->trRepo->add($t, true);

            return $this->redirectToRoute('dashboard_index');
        }

        return $this->render('transaction/add.html.twig', ['form' => $f]);
    }

    // TODO: Try to get this to be a delete method.
    #[Route('/delete/{id}', name: 'delete', methods: ['GET'])]
    public function delete(int $id): Response
    {
        $t = $this->trRepo->findOneBy(['id' => $id, 'user' => $this->getuser()]);
        if (!$t) {
            throw $this->createNotFoundException("No transaction with id $id.");
        }

        $this->trRepo->remove($t, true);

        return $this->redirectToRoute('transaction_index');
    }

    #[Route('/edit/{id}', name: 'edit', methods: ['GET', 'POST'])]
    public function edit(int $id, Request $request): Response
    {
        $t = $this->trRepo->findOneBy(['id' => $id, 'user' => $this->getuser()]);
        if (!$t) {
            throw $this->createNotFoundException("No transaction with id $id.");
        }

        $f = $this->createForm(TransactionType::class, $t);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $this->doctrine->getManager()->flush();

            return $this->redirectToRoute('transaction_index');
        }

        return $this->render('transaction/edit.html.twig', [
            'form' => $f,
            'transaction' => $t,
        ]);
    }

    #[Route('', name: 'index', methods: ['GET'])]
    public function index(Request $request): Response
    {
        $c = null;
        if ($cid = $request->query->get('category')) {
            $c = $this->catRepo->findOneBy(['id' => $cid, 'user' => $this->getUser()]);
        }
        if ($from = $request->query->get('from')) {
            try {
                $from = new DateTime($from);
            } catch (Exception) {
                $from = new DateTime();
            }
        }
        if ($to = $request->query->get('to')) {
            try {
                $to = new DateTime($to);
            } catch (Exception) {
                $to = new DateTime();
            }
        }

        $transactions = $this->trRepo->getFilteredTransactions($this->getUser(), $c, $from, $to);

        return $this->render('transaction/index.html.twig', [
            'categories' => $this->catRepo->findBy(['user' => $this->getUser()]),
            'filter' => [
                'categoryId' => $cid,
            ],
            'from' => $from,
            'stats' => $this->calculator->getTransactionsBalance($transactions),
            'to' => $to,
            'transactions' => $transactions,
        ]);
    }

    #[Route('/view/{id}', name: 'view', methods: ['GET'])]
    public function view(int $id): Response
    {
        $t = $this->trRepo->findOneBy(['id' => $id, 'user' => $this->getUser()]);
        if (!$t) {
            throw $this->createNotFoundException("No transaction with id $id.");
        }

        return $this->render('transaction/view.html.twig', ['transaction' => $t]);
    }
}
