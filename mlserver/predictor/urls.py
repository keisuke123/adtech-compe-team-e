from django.urls import path
from . import views

urlpatterns = [
    path('ctr', views.ctr)
]