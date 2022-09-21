<?php

namespace App\Controller\Private;

use App\Balance\Calculator;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/stats', name: 'stats_')]
class StatsController extends AbstractController
{
    #[Route('', name: 'index', methods: ['GET'])]
    public function index(Calculator $calculator): Response
    {
        return $this->render('stats/index.html.twig', [
            'stats' => $calculator->getAllTimeBalance($this->getUser()),
        ]);
    }
}
