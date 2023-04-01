<?php

namespace App\Form\Type;

use App\Entity\User;
use App\Form\Type\Element\SubmitIconButtonType;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\PasswordType;
use Symfony\Component\Form\Extension\Core\Type\RepeatedType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolver;
use Symfony\Component\Validator\Constraints\Length;
use Symfony\Component\Validator\Constraints\NotBlank;

class ResetPasswordType extends AbstractType
{
    public function buildForm(FormBuilderInterface $builder, array $options): void
    {
        $builder
            ->add('password', RepeatedType::class, [
                'attr' => ['autocomplete' => 'new-password'],
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
                'second_options' => [
                    'attr' => ['placeholder' => 'Repeat password'],
                    'label' => 'Repeat password:',
                ],
                // Instead of being set onto the object directly, this is read and hashed in the controller.
                'mapped' => false,
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
