<?php

namespace App\Form\Type;

use App\Entity\Category;
use App\Form\Type\Element\ResetIconButtonType;
use App\Form\Type\Element\SubmitIconButtonType;
use CapsuleType;
use Symfony\Component\Form\AbstractType;
use Symfony\Component\Form\Extension\Core\Type\ColorType;
use Symfony\Component\Form\Extension\Core\Type\FormType;
use Symfony\Component\Form\Extension\Core\Type\TextType;
use Symfony\Component\Form\FormBuilderInterface;
use Symfony\Component\OptionsResolver\OptionsResolver;

class CategoryType extends AbstractType
{
    public function buildForm(FormBuilderInterface $builder, array $options): void
    {
        $builder
            ->add('name', TextType::class, [
                'attr' => ['placeholder' => 'Name'],
                'label' => 'Name:',
                'row_attr' => ['class' => 'form-group'],
            ])
            ->add('color', ColorType::class, [
                'label' => 'Color:',
                'row_attr' => ['class' => 'form-group'],
            ])
            ->add('textColor', ColorType::class, [
                'label' => 'Text color:',
                'row_attr' => ['class' => 'form-group'],
            ])
            ->add('description', TextType::class, [
                'attr' => ['placeholder' => 'Description'],
                'empty_data' => '',
                'label' => 'Description:',
                'required' => false,
                'row_attr' => ['class' => 'form-group'],
            ])
            ->add(
                $builder->create('actions', FormType::class, [
                    'attr' => ['class' => 'form-actions'],
                    'inherit_data' => true,
                    'label' => false,
                ])
                ->add('save', SubmitIconButtonType::class, [
                    'attr' => ['class' => 'btn btn-icon'],
                ])
                ->add('cancel', ResetIconButtonType::class, [
                    'attr' => ['class' => 'btn btn-icon btn-cancel go-back'],
                ])
            );
    }

    public function configureOptions(OptionsResolver $resolver): void
    {
        $resolver->setDefaults(['data_class' => Category::class]);
    }
}
