import pytest

from pyinaturalist.formatters import (
    format_controlled_terms,
    format_identifications,
    format_observations,
    format_places,
    format_projects,
    format_species_counts,
    format_taxa,
    format_users,
    simplify_observations,
)
from test.conftest import load_sample_data

controlled_term_1 = load_sample_data('get_controlled_terms.json')['results'][0]
controlled_term_2 = load_sample_data('get_controlled_terms.json')['results'][1]
identification_1 = load_sample_data('get_identifications.json')['results'][0]
observation_1 = load_sample_data('get_observation.json')['results'][0]
observation_2 = load_sample_data('get_observations_node_page1.json')['results'][0]
place_1 = load_sample_data('get_places_by_id.json')['results'][0]
place_2 = load_sample_data('get_places_by_id.json')['results'][1]
project_1 = load_sample_data('get_projects.json')['results'][0]
project_2 = load_sample_data('get_projects.json')['results'][1]
species_count_1 = load_sample_data('get_observation_species_counts.json')['results'][0]
species_count_2 = load_sample_data('get_observation_species_counts.json')['results'][1]
taxon_1 = load_sample_data('get_taxa.json')['results'][0]
taxon_2 = load_sample_data('get_taxa.json')['results'][2]
user_1 = load_sample_data('get_users_autocomplete.json')['results'][0]
user_2 = load_sample_data('get_users_autocomplete.json')['results'][1]


def get_variations(response_object):
    """Formatting functions should accept any of these variations"""
    return [{'results': [response_object]}, [response_object], response_object]


@pytest.mark.parametrize('input', get_variations(controlled_term_1))
def test_format_controlled_terms(input):
    assert format_controlled_terms(input) == (
        '[      12] Plant Phenology\n'
        '    [      21] No Evidence of Flowering\n'
        '    [      13] Flowering\n'
        '    [      14] Fruiting\n'
        '    [      15] Flower Budding'
    )


def test_format_controlled_terms__align():
    assert format_controlled_terms([controlled_term_1, controlled_term_2], align=True) == (
        '[      12] Plant Phenology               \n'
        '    [      21] No Evidence of Flowering\n'
        '    [      13] Flowering\n'
        '    [      14] Fruiting\n'
        '    [      15] Flower Budding\n'
        '[       9] Sex                           \n'
        '    [      10] Female\n'
        '    [      11] Male\n'
        '    [      20] Cannot Be Determined'
    )


@pytest.mark.parametrize('input', get_variations(identification_1))
def test_format_identifications(input):
    assert format_identifications(input) == (
        '[155554373] [60132] Species: unknown '
        'identified by jkcook on 2021-02-18T20:31:32-06:00'
    )


@pytest.mark.parametrize('input', get_variations(observation_1))
def test_format_observation(input):
    assert format_observations(input) == (
        '[16227955] [493595] Species: Lixus bardanae '
        'observed by niconoe on 2018-09-05 at 54 rue des Badauds'
    )


def test_format_observation__align():
    expected_str = (
        '[16227955] [  493595]      Species                                          Lixus bardanae\n'
        '    observed by niconoe on 2018-09-05 at 54 rue des Badauds\n'
        '[57754375] [   48662]      Species                              Danaus plexippus (Monarch)\n'
        '    observed by samroom on 2020-08-27 at Railway Ave, Wilcox, SK, CA'
    )
    assert format_observations([observation_1, observation_2], align=True) == expected_str


@pytest.mark.parametrize('input', get_variations(species_count_1))
def test_format_species_counts(input):
    expected_str = '[48484] Species: Harmonia axyridis (Asian Lady Beetle): 31'
    assert format_species_counts(input) == expected_str


def test_format_species_counts__align():
    expected_str = (
        '[   48484]      Species                   Harmonia axyridis (Asian Lady Beetle): 31\n'
        '[   51702]      Species   Coccinella septempunctata (Seven-spotted Lady Beetle): 19'
    )
    assert format_species_counts([species_count_1, species_count_2], align=True) == expected_str


@pytest.mark.parametrize('input', get_variations(taxon_1))
def test_format_taxa__with_common_name(input):
    expected_str = '[70118] Species: Nicrophorus vespilloides (Lesser Vespillo Burying Beetle)'
    assert format_taxa(input) == expected_str


@pytest.mark.parametrize('input', get_variations(taxon_2))
def test_format_taxon__without_common_name(input):
    assert format_taxa(input) == '[124162] Species: Temnostoma vespiforme'


def test_format_taxa__align():
    expected_str = (
        '[   70118]      Species Nicrophorus vespilloides (Lesser Vespillo Burying Beetle)\n'
        '[  124162]      Species                                   Temnostoma vespiforme'
    )
    x = format_taxa([taxon_1, taxon_2], align=True)
    print(x)
    assert format_taxa([taxon_1, taxon_2], align=True) == expected_str


def test_format_taxon__invalid():
    assert format_taxa(None) == 'unknown taxon'


def test_simplify_observation():
    simplified_obs = simplify_observations(observation_1)
    # Not much worth testing here, just make sure it returns something that can be formatted
    assert format_observations(simplified_obs)


@pytest.mark.parametrize('input', get_variations(place_1))
def test_format_places(input):
    assert format_places(input) == '[93735] Springbok, Northern Cape'


def test_format_places__align():
    assert format_places([place_1, place_2], align=True) == (
        '[   93735] Springbok, Northern Cape\n'
        '[   89191] Conservation Area Riversdale, CA, NZ'
    )


@pytest.mark.parametrize('input', get_variations(project_1))
def test_format_projects(input):
    assert format_projects(input) == '[8291] PNW Invasive Plant EDDR'


def test_format_projects__align():
    assert format_projects([project_1, project_2], align=True) == (
        '[    8291] PNW Invasive Plant EDDR\n'
        '[   19200] King County (WA) Noxious and Invasive Weeds'
    )


@pytest.mark.parametrize('input', get_variations(user_1))
def test_format_users(input):
    assert format_users(input) == '[886482] niconoe (Nicolas Noé)'


def test_format_users__align():
    assert format_users([user_1, user_2], align=True) == (
        '[  886482] niconoe (Nicolas Noé)\n' '[ 2909130] niconoerbo'
    )
