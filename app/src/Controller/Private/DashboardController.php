<?php

namespace App\Controller\Private;

use App\Entity\TransactionTemplate;
use Doctrine\Persistence\ManagerRegistry;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

class DashboardController extends AbstractController
{
    #[Route('', methods: ['GET'], name: 'dashboard_index')]
    public function index(ManagerRegistry $doctrine): Response
    {
        return $this->render('dashboard/index.html.twig', [
            'transactionTemplates' => $doctrine->getRepository(TransactionTemplate::class)->findBy(
                [],
                ['position' => 'ASC']
            )
        ]);
    }
}
