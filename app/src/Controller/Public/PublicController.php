<?php

namespace App\Controller\Public;

use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

class PublicController extends AbstractController
{
    #[Route('/', methods: ['GET'], name: 'public_home')]
    public function home(): Response
    {
        return $this->render('public/home.html.twig');
    }
}
