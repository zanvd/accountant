<?php

namespace App\Controller\Private;

use App\Balance\Calculator;
use App\Stats\Category;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/stats', name: 'stats_')]
class StatsController extends AbstractController
{
    private Category $catStats;

    public function __construct(Category $catStats)
    {
        $this->catStats = $catStats;
    }

    #[Route('', name: 'index', methods: ['GET'])]
    public function index(Calculator $calculator): Response
    {
        return $this->render('stats/index.html.twig', [
            'catStats' => $this->catStats->calcAll(),
            'stats' => $calculator->getAllTimeBalance($this->getUser()),
        ]);
    }
}
