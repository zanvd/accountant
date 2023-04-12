<?php

namespace App\Form\Type;

use App\Entity\User;
use App\Enum\UserValidationType;
use App\Form\Type\Element\SubmitIconButtonType;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\EmailType;
use Symfony\Component\Form\Extension\Core\Type\PasswordType;
use Symfony\Component\Form\Extension\Core\Type\RepeatedType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolver;

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
                'first_options' => [
                    'attr' => [
                        'autocomplete' => 'new-password',
                        'placeholder' => 'Password',
                    ],
                    'label' => 'Password:',
                ],
                'second_options' => [
                    'attr' => [
                        'autocomplete' => 'new-password',
                        'placeholder' => 'Repeat password',
                    ],
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
        $resolver->setDefaults([
            'data_class' => User::class,
            'validation_groups' => [UserValidationType::Register->value],
        ]);
    }
}
