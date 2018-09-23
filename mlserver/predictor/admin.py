from django.contrib import admin
from solo.admin import SingletonModelAdmin

# Register your models here.
from predictor.models import CTRPredictor, Preprocessor

admin.site.register(CTRPredictor, SingletonModelAdmin)
admin.site.register(Preprocessor, SingletonModelAdmin)

preprocessor = Preprocessor.get_solo()
ctr_predictor = CTRPredictor.get_solo()
