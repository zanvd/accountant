<?php

namespace App\Controller\Private;

use App\Entity\Category;
use App\Form\Type\CategoryType;
use App\Repository\CategoryRepository;
use Doctrine\Persistence\ManagerRegistry;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\Routing\Annotation\Route;

#[Route('/category', name: 'category_')]
class CategoryController extends AbstractController
{
    private CategoryRepository $catRepo;
    private ManagerRegistry $doctrine;

    public function __construct(ManagerRegistry $doctrine)
    {
        $this->catRepo = $doctrine->getRepository(Category::class);
        $this->doctrine = $doctrine;
    }

    #[Route('/add', methods: ['GET', 'POST'], name: 'add')]
    public function add(Request $request): Response
    {
        $c = new Category();
        $c->setUser($this->getUser());
        $f = $this->createForm(CategoryType::class, $c);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $this->catRepo->add($c, true);

            return $this->redirectToRoute('category_index');
        }

        return $this->renderForm('category/add.html.twig', ['form' => $f]);
    }

    // TODO: Try to get this to be a delete method.
    #[Route('/delete/{id}', methods: ['GET'], name: 'delete')]
    public function delete(int $id): Response
    {
        $c = $this->catRepo->findOneBy(['id' => $id, 'user' => $this->getuser()]);
        if (!$c) throw $this->createNotFoundException("No category with id $id.");

        $this->catRepo->remove($c, true);

        return $this->redirectToRoute('category_index');
    }

    #[Route('/edit/{id}', methods: ['GET', 'POST'], name: 'edit')]
    public function edit(int $id, Request $request): Response
    {
        $c = $this->catRepo->findOneBy(['id' => $id, 'user' => $this->getuser()]);
        if (!$c) throw $this->createNotFoundException("No category with id $id.");

        $f = $this->createForm(CategoryType::class, $c);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $this->doctrine->getManager()->flush();

            return $this->redirectToRoute('category_index');
        }

        return $this->renderForm('category/edit.html.twig', [
            'category' => $c,
            'form' => $f,
        ]);
    }

    #[Route('', methods: ['GET'], name: 'index')]
    public function index(): Response
    {
        return $this->render('category/index.html.twig', [
            'categories' => $this->catRepo->findBy(['user' => $this->getUser()])
        ]);
    }

    #[Route('/view/{id}', methods: ['GET'], name: 'view')]
    public function view(int $id): Response
    {
        $c = $this->catRepo->findOneBy(['id' => $id, 'user' => $this->getUser()]);
        if (!$c) throw $this->createNotFoundException("No category with id $id.");

        return $this->render('category/view.html.twig', ['category' => $c]);
    }
}
