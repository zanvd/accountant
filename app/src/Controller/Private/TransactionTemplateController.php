<?php

namespace App\Controller\Private;

use App\Entity\TransactionTemplate;
use App\Form\Type\TransactionTemplateType;
use App\Repository\TransactionTemplateRepository;
use Doctrine\Persistence\ManagerRegistry;
use Doctrine\Persistence\ObjectRepository;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/transaction-template', name: 'transaction_template_')]
class TransactionTemplateController extends AbstractController
{
    private ManagerRegistry $doctrine;
    private ObjectRepository|TransactionTemplateRepository $trTeRepo;

    public function __construct(ManagerRegistry $doctrine)
    {
        $this->doctrine = $doctrine;
        $this->trTeRepo = $doctrine->getRepository(TransactionTemplate::class);
    }

    #[Route('/add', name: 'add', methods: ['GET', 'POST'])]
    public function add(Request $request): Response
    {
        $tt = new TransactionTemplate();
        $tt->setUser($this->getUser());
        $f = $this->createForm(TransactionTemplateType::class, $tt);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $this->trTeRepo->add($tt, true);

            return $this->redirectToRoute('transaction_template_index');
        }

        return $this->renderForm('transactionTemplate/add.html.twig', ['form' => $f]);
    }

    // TODO: Try to get this to be a delete method.
    #[Route('/delete/{id}', name: 'delete', methods: ['GET'])]
    public function delete(int $id): Response
    {
        $tt = $this->trTeRepo->findOneBy(['id' => $id, 'user' => $this->getuser()]);
        if (!$tt) throw $this->createNotFoundException("No transaction template with id $id.");

        $this->trTeRepo->remove($tt, true);

        return $this->redirectToRoute('transaction_template_index');
    }

    #[Route('/edit/{id}', name: 'edit', methods: ['GET', 'POST'])]
    public function edit(int $id, Request $request): Response
    {
        $tt = $this->trTeRepo->findOneBy(['id' => $id, 'user' => $this->getuser()]);
        if (!$tt) throw $this->createNotFoundException("No transaction template with id $id.");

        $f = $this->createForm(TransactionTemplateType::class, $tt);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $this->doctrine->getManager()->flush();

            return $this->redirectToRoute('transaction_template_index');
        }

        return $this->renderForm('transactionTemplate/edit.html.twig', [
            'form' => $f,
            'transactionTemplate' => $tt,
        ]);
    }

    #[Route('', name: 'index', methods: ['GET'])]
    public function index(): Response
    {
        return $this->render('transactionTemplate/index.html.twig', [
            'transactionTemplates' => $this->trTeRepo->findBy(['user' => $this->getUser()]),
        ]);
    }

    #[Route('/view/{id}', name: 'view', methods: ['GET'])]
    public function view(int $id): Response
    {
        $tt = $this->trTeRepo->findOneBy(['id' => $id, 'user' => $this->getUser()]);
        if (!$tt) throw $this->createNotFoundException("No transaction template with id $id.");

        return $this->render('transactionTemplate/view.html.twig', ['transactionTemplate' => $tt]);
    }
}
