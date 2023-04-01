<?php

namespace App\Form\Type;

use App\Entity\User;
use App\Form\Type\Element\SubmitIconButtonType;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\EmailType;
use Symfony\Component\Form\Extension\Core\Type\PasswordType;
use Symfony\Component\Form\Extension\Core\Type\RepeatedType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolver;
use Symfony\Component\Validator\Constraints\Length;
use Symfony\Component\Validator\Constraints\NotBlank;

class RegistrationType extends AbstractType
{
    public function buildForm(FormBuilderInterface $builder, array $options): void
    {
        $builder
            ->add('username', TextType::class, [
                'attr' => [
                    'autocomplete' => 'off',
                    'placeholder' => 'Username',
                ],
                'label' => 'Username:',
            ])
            ->add('email', EmailType::class, [
                'attr' => [
                    'autocomplete' => 'off',
                    'placeholder' => 'Email',
                ],
                'label' => 'Email:',
            ])
            ->add('password', RepeatedType::class, [
                'constraints' => [
                    new NotBlank([
                        'message' => 'Please enter a password',
                    ]),
                    new Length([
                        // Max length allowed by Symfony for security reasons.
                        'max' => 4096,
                        'min' => 6,
                        'minMessage' => 'Your password should be at least {{ limit }} characters',
                    ]),
                ],
                'first_options' => [
                    'attr' => ['placeholder' => 'Password'],
                    'label' => 'Password:',
                ],
                // Instead of being set onto the object directly, this is read and hashed in the controller.
                'mapped' => false,
                'options' => [
                    'attr' => ['autocomplete' => 'new-password'],
                ],
                'second_options' => [
                    'attr' => ['placeholder' => 'Repeat password'],
                    'label' => 'Repeat password:',
                ],
                'type' => PasswordType::class,
            ])
            ->add('save', SubmitIconButtonType::class, [
                'attr' => ['class' => 'btn btn-icon'],
                'row_attr' => ['class' => 'form-actions'],
            ]);
    }

    public function configureOptions(OptionsResolver $resolver): void
    {
        $resolver->setDefaults(['data_class' => User::class]);
    }
}
