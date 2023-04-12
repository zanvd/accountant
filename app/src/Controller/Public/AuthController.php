<?php

namespace App\Controller\Public;

use App\Entity\User;
use App\Form\Type\ForgotPasswordType;
use App\Form\Type\LoginType;
use App\Form\Type\RegistrationType;
use App\Form\Type\ResetPasswordType;
use App\Repository\UserRepository;
use App\Security\EmailVerifier;
use App\Security\ResetPassword;
use Doctrine\Persistence\ManagerRegistry;
use Doctrine\Persistence\ObjectRepository;
use Exception;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\PasswordHasher\Hasher\UserPasswordHasherInterface;
use Symfony\Component\Routing\Annotation\Route;
use Symfony\Component\Security\Http\Authentication\AuthenticationUtils;
use SymfonyCasts\Bundle\VerifyEmail\Exception\VerifyEmailExceptionInterface;

#[Route('/', name: 'auth_')]
class AuthController extends AbstractController
{
    private EmailVerifier $emailVerifier;
    private ObjectRepository|UserRepository $userRepository;

    public function __construct(ManagerRegistry $doctrine, EmailVerifier $emailVerifier)
    {
        $this->emailVerifier = $emailVerifier;
        $this->userRepository = $doctrine->getRepository(User::class);
    }

    #[Route('/forgot-password', name: 'forgot_password', methods: ['GET', 'POST'])]
    public function forgotPassword(Request $request, ResetPassword $resetPassword): Response
    {
        $f = $this->createForm(ForgotPasswordType::class);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $u = $this->userRepository->findOneBy([
                'username' => $request->request->all($f->getName())['username'] ?? ''
            ]);
            if ($u && $resetPassword->initPasswordReset($u)) {
                $this->addFlash('success', 'We\'ve sent a password reset link to your email.');

                return $this->redirectToRoute('auth_forgot_password');
            } else {
                $this->addFlash('error', 'Could not send a password reset link to your email. Please, try again.');
            }
        }

        return $this->render('auth/forgot-password.html.twig', ['form' => $f]);
    }

    #[Route('/login', name: 'login', methods: ['GET', 'POST'])]
    public function login(AuthenticationUtils $authUtils): Response
    {
        $u = new User();
        $f = $this->createForm(LoginType::class, $u);

        return $this->render('auth/login.html.twig', [
            'error' => $authUtils->getLastAuthenticationError(),
            'form' => $f,
        ]);
    }

    #[Route('/logout', name: 'logout', methods: ['GET'])]
    public function logout(): void
    {
        throw new Exception('Should not be called.');
    }

    #[Route('/register', name: 'register', methods: ['GET', 'POST'])]
    public function register(UserPasswordHasherInterface $passHasher, Request $request): Response
    {
        $f = $this->createForm(RegistrationType::class);
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $u = $f->getData();
            $u->setPassword($passHasher->hashPassword($u, $f->get('password')->getData()));
            $this->userRepository->add($u, true);

            $this->emailVerifier->sendEmailConfirmation($u, 'auth_verify_email');

            return $this->redirectToRoute('dashboard_index');
        }

        return $this->render('auth/register.html.twig', ['form' => $f]);
    }

    #[Route('/reset-password', name: 'reset_password', methods: ['GET', 'POST'])]
    public function resetPassword(
        UserPasswordHasherInterface $passHasher,
        Request $request,
        ResetPassword $resetPassword
    ): Response {
        $t = $request->query->get('token', '');
        $un = $request->query->get('username', '');
        $cun = $resetPassword->getCachedData($t);
        if (!$cun || $un !== $cun) {
            $this->addFlash('error', 'Invalid token provided. Please, request a new link.');

            $this->redirectToRoute('auth_forgot_password');
        }

        $u = $this->userRepository->findOneBy(['username' => $un]);
        if (!$u) {
            $this->addFlash('error', 'Invalid username provided. Please, request a new link.');

            return $this->redirectToRoute('auth_forgot_password');
        }

        $f = $this->createForm(
            ResetPasswordType::class,
            options: [
                'action' => $this->generateUrl(
                    'auth_reset_password',
                    [
                        'token' => $t,
                        'username' => $un,
                    ]
                ),
            ]
        );
        $f->handleRequest($request);
        if ($f->isSubmitted() && $f->isValid()) {
            $u->setPassword($passHasher->hashPassword($u, $f->get('password')->getData()));
            $this->userRepository->add($u, true);

            $resetPassword->clearCachedData($t);

            return $this->redirectToRoute('auth_login');
        }

        return $this->render('auth/reset-password.html.twig', ['form' => $f]);
    }

    #[Route('/verify-email', name: 'verify_email', methods: ['GET'])]
    public function verifyUserEmail(Request $request): Response
    {
        $id = $request->get('id');
        if (is_null($id)) {
            return $this->redirectToRoute('auth_register');
        }

        $u = $this->userRepository->find($id);
        if (!$u) {
            return $this->redirectToRoute('auth_register');
        }

        try {
            $this->emailVerifier->handleEmailConfirmation($request, $u);
        } catch (VerifyEmailExceptionInterface $e) {
            $this->addFlash('verify_email_error', $e->getReason());

            return $this->redirectToRoute('auth_register');
        }

        $this->addFlash('success', 'Your email address has been verified.');

        return $this->redirectToRoute('auth_login');
    }
}
