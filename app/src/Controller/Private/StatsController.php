<?php

namespace App\Controller\Private;

use App\Repository\TransactionRepository;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/stats', name: 'stats_')]
class StatsController extends AbstractController
{
    #[Route('', name: 'index', methods: ['GET'])]
    public function index(TransactionRepository $transactionRepository): Response
    {
        $i = 0;
        $o = 0;
        foreach ($transactionRepository->findBy(['user' => $this->getUser()]) as $t) {
            if ($t->getAmount() > 0)
                $i += $t->getAmount();
            else
                $o += $t->getAmount();
        }

        return $this->render('stats/index.html.twig', [
            'stats' => [
                'income' => $i,
                'outcome' => $o,
                'savings' => $i + $o,
            ]
        ]);
    }
}
