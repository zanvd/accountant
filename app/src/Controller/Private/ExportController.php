<?php

namespace App\Controller\Private;

use App\Exporter\ExporterInterface;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/export', name: 'export_')]
class ExportController extends AbstractController
{
    private ExporterInterface $exporter;

    public function __construct(ExporterInterface $exporter)
    {
        $this->exporter = $exporter;
    }

    #[Route('/add', name: 'add', methods: ['GET'])]
    public function add(): Response
    {
        return new Response($this->exporter->export(), headers: [
            'Content-Disposition' => 'filename="accountant-export-' . date('m-d-Y') . '.csv"',
            'Content-Type' => 'text/csv'
        ]);
    }

    #[Route('', name: 'index', methods: ['GET'])]
    public function index(): Response
    {
        return $this->render('export/index.html.twig');
    }
}
