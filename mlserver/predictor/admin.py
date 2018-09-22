from django.contrib import admin
from solo.admin import SingletonModelAdmin

# Register your models here.
from predictor.models import CTRPredictor

admin.site.register(CTRPredictor, SingletonModelAdmin)

ctr_predictor = CTRPredictor.get_solo()
